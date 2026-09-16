async function loadStats() {
  const res = await fetch('/api/stats/overview');
  const body = await res.json();
  const data = body.data || {};
  const container = document.getElementById('stats-cards');
  const cards = [
    { label: '证书总数', value: data.total_certificates || 0 },
    { label: '活跃证书', value: data.active_certificates || 0 },
    { label: '已过期', value: data.expired_certificates || 0 },
    { label: '已吊销', value: data.revoked_certificates || 0 },
    { label: '分类数', value: data.total_categories || 0 },
    { label: '持有人', value: data.total_holders || 0 },
    { label: '发证机构', value: data.total_issuers || 0 },
    { label: '待审续期', value: data.pending_renewals || 0 },
  ];
  container.innerHTML = '';
  cards.forEach(c => {
    const div = document.createElement('div');
    div.className = 'stat-card';
    div.innerHTML = `<div class="stat-value">${c.value}</div><div class="stat-label">${c.label}</div>`;
    container.appendChild(div);
  });
}

async function loadCertificates() {
  const res = await fetch('/api/certificates');
  const body = await res.json();
  const items = (body.data && body.data.items) || [];
  const tbody = document.querySelector('#cert-table tbody');
  tbody.innerHTML = '';
  items.forEach(c => {
    const tr = document.createElement('tr');
    const statusClass = 'status-' + c.status;
    tr.innerHTML = `
      <td>${c.code}</td>
      <td>${c.name}</td>
      <td>${c.category_id}</td>
      <td>${c.holder_id}</td>
      <td>${c.issuer_id}</td>
      <td>${c.issue_date}</td>
      <td>${c.expire_date}</td>
      <td><span class="status-badge ${statusClass}">${c.status}</span></td>
    `;
    tbody.appendChild(tr);
  });
}

async function loadExpiringSoon() {
  const res = await fetch('/api/stats/expiring-soon?limit=10');
  const body = await res.json();
  const items = body.data || [];
  const tbody = document.querySelector('#expiring-table tbody');
  tbody.innerHTML = '';
  items.forEach(c => {
    const tr = document.createElement('tr');
    tr.innerHTML = `<td>${c.code}</td><td>${c.name}</td><td>${c.expire_date}</td><td>${c.days_left}</td>`;
    tbody.appendChild(tr);
  });
}

async function init() {
  await loadStats();
  await loadCertificates();
  await loadExpiringSoon();
}

init();
