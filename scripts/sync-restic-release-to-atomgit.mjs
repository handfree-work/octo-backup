import { mkdtemp, readFile, rm, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { fetch, ProxyAgent } from 'undici';

const githubRepository = process.env.GITHUB_REPOSITORY ?? 'restic/restic';
const atomgitRepository = process.env.ATOMGIT_REPOSITORY ?? 'handfree-work/restic';
const apiBase = process.env.ATOMGIT_API_BASE ?? 'https://api.atomgit.com/api/v5';
const proxy = process.env.HTTPS_PROXY ?? 'http://127.0.0.1:10811';
const tagArg = process.argv[2];
const atomgitToken = process.env.ATOMGIT_TOKEN;
const githubToken = process.env.GH_TOKEN ?? process.env.GITHUB_TOKEN;

if (!atomgitToken) throw new Error('缺少环境变量 ATOMGIT_TOKEN');
if (!/^[^/]+\/[^/]+$/.test(atomgitRepository)) throw new Error('AtomGit 仓库必须使用 owner/repo 格式');

const githubDispatcher = new ProxyAgent(proxy);
const githubHeaders = { accept: 'application/vnd.github+json' };
if (githubToken) githubHeaders.authorization = `Bearer ${githubToken}`;

async function githubJson(url) {
  const response = await fetch(url, { headers: githubHeaders, dispatcher: githubDispatcher });
  if (!response.ok) throw new Error(`GitHub API ${response.status}: ${await response.text()}`);
  return response.json();
}

async function atomgitJson(url, options = {}) {
  const response = await fetch(url, { ...options, headers: { accept: 'application/json', ...(options.headers ?? {}) } });
  if (!response.ok) throw new Error(`AtomGit API ${response.status}: ${await response.text()}`);
  return response.json();
}

const releaseUrl = tagArg
  ? `https://api.github.com/repos/${githubRepository}/releases/tags/${tagArg.startsWith('v') ? tagArg : `v${tagArg}`}`
  : `https://api.github.com/repos/${githubRepository}/releases/latest`;
const release = await githubJson(releaseUrl);
const tag = release.tag_name;
console.log(`GitHub Release 版本: ${tag}`);
if (!/^v\d+\.\d+\.\d+(?:[-.][0-9A-Za-z.-]+)?$/.test(tag)) throw new Error(`GitHub Release 标签格式无效: ${tag}`);
if (!release.assets?.length) throw new Error(`GitHub Release ${tag} 没有可同步的安装包`);
const windowsZipAssets = release.assets.filter((asset) => /_windows_[^/]+\.zip$/i.test(asset.name));
if (!windowsZipAssets.length) throw new Error(`GitHub Release ${tag} 缺少 Windows ZIP 安装包`);

const directory = await mkdtemp(join(tmpdir(), 'restic-release-'));
try {
  const [owner, repo] = atomgitRepository.split('/');
  const base = `${apiBase}/repos/${owner}/${repo}/releases`;
  const token = encodeURIComponent(atomgitToken);
  const body = release.body?.trim() || `Release ${tag}`;
  try {
    await atomgitJson(`${base}?access_token=${token}`, { method: 'POST', headers: { 'content-type': 'application/json' }, body: JSON.stringify({ name: release.name || tag, tag_name: tag, body, release_status: 'latest' }) });
  } catch (error) {
    if (!String(error).includes('409')) throw error;
    console.log(`AtomGit release ${tag} 已存在，继续同步资产`);
  }

  const atomgitRelease = await atomgitJson(`${base}/${encodeURIComponent(tag)}?access_token=${token}`);
  const existingAssets = new Set((atomgitRelease.assets ?? []).map((asset) => asset.name));

  for (const asset of release.assets) {
    if (existingAssets.has(asset.name)) {
      console.log(`跳过已存在的 ${asset.name}`);
      continue;
    }
    console.log(`下载 ${asset.name}`);
    const response = await fetch(asset.browser_download_url, { headers: githubHeaders, dispatcher: githubDispatcher });
    if (!response.ok) throw new Error(`下载 ${asset.name} 失败: ${response.status}`);
    await writeFile(join(directory, asset.name), Buffer.from(await response.arrayBuffer()));
    const upload = await atomgitJson(`${base}/${tag}/upload_url?access_token=${token}&file_name=${encodeURIComponent(asset.name)}`);
    if (!upload.url) throw new Error(`未获取到 ${asset.name} 的 AtomGit 上传地址`);
    const headers = upload.headers ?? {};
    if (!Object.keys(headers).some((key) => key.toLowerCase() === 'content-type')) {
      headers['content-type'] = asset.name.toLowerCase().endsWith('.zip') ? 'application/zip' : 'application/octet-stream';
    }
    console.log(`上传 ${asset.name}`);
    const result = await fetch(upload.url, { method: 'PUT', headers, body: await readFile(join(directory, asset.name)) });
    const text = await result.text();
    if (!result.ok || /Fail to call origin|"success"\s*:\s*false/.test(text)) {
      throw new Error(`AtomGit 上传失败 (${asset.name}): ${text.replaceAll(atomgitToken, '[REDACTED]')}`);
    }
    existingAssets.add(asset.name);
  }
  console.log(`已同步 ${githubRepository} ${tag} 到 ${atomgitRepository}`);
} finally {
  await rm(directory, { recursive: true, force: true });
}
