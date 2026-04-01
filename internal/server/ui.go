package server

import "net/http"

func (s *Server) handleUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Ticker — Stockyard</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<style>*{margin:0;padding:0;box-sizing:border-box}body{background:#1a1410;color:#f0e6d3;font-family:'JetBrains Mono',monospace;padding:2rem}
.hdr{font-size:.7rem;color:#a0845c;letter-spacing:3px;text-transform:uppercase;margin-bottom:2rem;border-bottom:2px solid #8b3d1a;padding-bottom:.8rem}
.cards{display:grid;grid-template-columns:repeat(3,1fr);gap:1rem;margin-bottom:2rem}.card{background:#241e18;border:1px solid #2e261e;padding:1rem}.card-val{font-size:1.6rem;font-weight:700;display:block}.card-lbl{font-size:.55rem;letter-spacing:2px;text-transform:uppercase;color:#a0845c;margin-top:.2rem}
table{width:100%;border-collapse:collapse;font-size:.72rem;margin-top:1rem}th{background:#2e261e;padding:.4rem .6rem;text-align:left;color:#c4a87a;font-size:.6rem;letter-spacing:1px;text-transform:uppercase}td{padding:.4rem .6rem;border-bottom:1px solid #2e261e;color:#bfb5a3}.empty{color:#7a7060;text-align:center;padding:2rem;font-style:italic}
.section{margin-bottom:2rem}.section h2{font-size:.65rem;letter-spacing:3px;text-transform:uppercase;color:#e8753a;margin-bottom:.5rem}
.pos{color:#5ba86e}.neg{color:#c0392b}
</style></head><body>
<div class="hdr">Stockyard · Ticker</div>
<div class="cards"><div class="card"><span class="card-val" id="s-accts">—</span><span class="card-lbl">Accounts</span></div><div class="card"><span class="card-val" id="s-txns">—</span><span class="card-lbl">Transactions</span></div><div class="card"><span class="card-val" id="s-bal">—</span><span class="card-lbl">Net Balance</span></div></div>
<div class="section"><h2>Recent Transactions</h2>
<table><thead><tr><th>Date</th><th>Description</th><th>Category</th><th style="text-align:right">Amount</th></tr></thead><tbody id="txn-body"></tbody></table></div>
<script>
async function refresh(){
  try{const s=await(await fetch('/api/status')).json();document.getElementById('s-accts').textContent=s.accounts||0;document.getElementById('s-txns').textContent=s.transactions||0;document.getElementById('s-bal').textContent='$'+((s.net_balance_cents||0)/100).toFixed(2);}catch(e){}
  try{const d=await(await fetch('/api/transactions')).json();const ts=d.transactions||[];const tb=document.getElementById('txn-body');
  tb.innerHTML=ts.length?ts.map(t=>{const cls=t.amount_cents>=0?'pos':'neg';return '<tr><td>'+t.date+'</td><td>'+esc(t.description)+'</td><td>'+esc(t.category||'—')+'</td><td style="text-align:right" class="'+cls+'">$'+(t.amount_cents/100).toFixed(2)+'</td></tr>';}).join(''):'<tr><td colspan="4" class="empty">No transactions</td></tr>';}catch(e){}
}
function esc(s){const d=document.createElement('div');d.textContent=s||'';return d.innerHTML;}
refresh();setInterval(refresh,8000);
</script></body></html>`))
}
