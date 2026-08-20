async function fetchSnapshot() {
  const res = await fetch('/v1/blocks');
  if (!res.ok) throw new Error('snapshot failed');
  return res.json();
}

function renderGrid(data) {
  const grid = document.getElementById('grid');
  const summary = document.getElementById('summary');
  grid.innerHTML = '';
  summary.textContent = data.summary || '无数据';

  const blocks = data.blocks || [];
  if (blocks.length === 0) {
    grid.innerHTML = '<p class="summary">暂无分区上报</p>';
    return;
  }

  blocks.forEach(b => {
    const div = document.createElement('div');
    div.className = 'cell';
    if (b.conflict) div.classList.add('conflict');
    else if (b.occupied) div.classList.add('occupied');
    else div.classList.add('free');

    const sources = (b.sources || []).filter(s => s.occupied).map(s => s.source).join(', ');
    div.innerHTML = '<div class="id">' + b.block_id + '</div>' +
      '<div class="state">' + (b.conflict ? '冲突' : (b.occupied ? '占用' : '空闲')) + '</div>' +
      '<div class="state">' + (sources || '-') + '</div>';
    grid.appendChild(div);
  });
}

function renderConflicts(data) {
  const ul = document.getElementById('conflicts');
  ul.innerHTML = '';
  const records = data.conflicts || [];
  if (records.length === 0) {
    ul.innerHTML = '<li>无冲突</li>';
    return;
  }
  records.forEach(r => {
    const li = document.createElement('li');
    li.textContent = '分区 ' + r.block_id + ': ' + r.message + (r.sources ? ' [' + r.sources.join(', ') + ']' : '');
    ul.appendChild(li);
  });
}

async function refresh() {
  try {
    const data = await fetchSnapshot();
    renderGrid(data);
    renderConflicts(data);
  } catch (e) {
    document.getElementById('summary').textContent = '加载失败: ' + e.message;
  }
}

document.getElementById('refreshBtn').addEventListener('click', refresh);

document.getElementById('clearanceForm').addEventListener('submit', async (ev) => {
  ev.preventDefault();
  const raw = document.getElementById('routeBlocks').value.trim();
  const blocks = raw.split(',').map(s => parseInt(s.trim(), 10)).filter(n => !isNaN(n));
  const res = await fetch('/v1/clearance/check', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ route_id: 'ui', blocks })
  });
  const data = await res.json();
  document.getElementById('clearanceResult').textContent = JSON.stringify(data, null, 2);
});

document.getElementById('ingestForm').addEventListener('submit', async (ev) => {
  ev.preventDefault();
  const source = document.getElementById('sourceId').value.trim();
  const blockId = parseInt(document.getElementById('blockId').value, 10);
  const occupied = parseInt(document.getElementById('occupied').value, 10);
  const seq = Math.floor(Date.now() / 1000) >>> 0;

  const res = await fetch('/v1/frames/encode', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ block_id: blockId, occupied: occupied, seq: seq })
  });
  if (!res.ok) {
    document.getElementById('ingestResult').textContent = 'encode failed';
    return;
  }
  const enc = await res.json();
  const frameRes = await fetch('/v1/frames', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/octet-stream',
      'X-Source-ID': source
    },
    body: hexToBytes(enc.hex)
  });
  const data = await frameRes.json();
  document.getElementById('ingestResult').textContent = JSON.stringify(data, null, 2);
  await refresh();
});

function hexToBytes(hex) {
  const out = new Uint8Array(hex.length / 2);
  for (let i = 0; i < hex.length; i += 2) {
    out[i / 2] = parseInt(hex.substr(i, 2), 16);
  }
  return out;
}

refresh();
setInterval(refresh, 5000);
