package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Mainspring</title>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital,wght@0,400;0,700;1,400&family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center;font-family:var(--mono)}
.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center}
.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.job{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.job:hover{border-color:var(--leather)}
.job.disabled{opacity:.5}
.job-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.job-name{font-size:.88rem;font-weight:700}
.job-schedule{font-family:var(--mono);font-size:.65rem;color:var(--gold);margin-top:.1rem}
.job-cmd{font-family:var(--mono);font-size:.6rem;color:var(--cm);margin-top:.3rem;background:var(--bg);padding:.3rem .5rem;border:1px solid var(--bg3);word-break:break-all}
.job-meta{font-family:var(--mono);font-size:.55rem;color:var(--cm);margin-top:.35rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.job-actions{display:flex;gap:.3rem;align-items:center;flex-shrink:0}
.result-ok{color:var(--green)}.result-fail{color:var(--red)}
.toggle{position:relative;display:inline-block;width:32px;height:18px}
.toggle input{opacity:0;width:0;height:0}
.sl{position:absolute;cursor:pointer;inset:0;background:var(--bg3);transition:.2s;border-radius:18px}
.sl:before{content:'';position:absolute;height:14px;width:14px;left:2px;bottom:2px;background:var(--cm);transition:.2s;border-radius:50%}
.toggle input:checked+.sl{background:var(--green)}
.toggle input:checked+.sl:before{transform:translateX(14px);background:var(--cream)}
.btn{font-family:var(--mono);font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.btn-run{border-color:var(--green);color:var(--green)}.btn-run:hover{background:var(--green);color:#fff}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:480px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.fr .hint{font-size:.5rem;color:var(--cm);margin-top:.15rem}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}.row2{grid-template-columns:1fr}.toolbar{flex-direction:column}.search{width:100%}}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> MAINSPRING</h1><button class="btn btn-p" onclick="openForm()">+ New Job</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar"><input class="search" id="search" placeholder="Search jobs..." oninput="render()"></div>
<div id="list"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',jobs=[],editId=null;

async function load(){var r=await fetch(A+'/jobs').then(function(r){return r.json()});jobs=r.jobs||[];renderStats();render();}

function renderStats(){
var total=jobs.length;
var active=jobs.filter(function(j){return j.enabled}).length;
var totalRuns=jobs.reduce(function(s,j){return s+j.run_count},0);
var totalFails=jobs.reduce(function(s,j){return s+j.fail_count},0);
document.getElementById('stats').innerHTML=[
{l:'Jobs',v:total},{l:'Active',v:active,c:active>0?'var(--green)':''},{l:'Total Runs',v:totalRuns},{l:'Failures',v:totalFails,c:totalFails>0?'var(--red)':''}
].map(function(x){return '<div class="st"><div class="st-v" style="'+(x.c?'color:'+x.c:'')+'">'+x.v+'</div><div class="st-l">'+x.l+'</div></div>'}).join('');
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var f=jobs;
if(q)f=f.filter(function(j){return(j.name||'').toLowerCase().includes(q)||(j.command||'').toLowerCase().includes(q)||(j.schedule||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No scheduled jobs. Create one to get started.</div>';return;}
var h='';f.forEach(function(j){
var ok=j.last_result==='success'||j.last_result==='ok';
var fail=j.last_result==='failure'||j.last_result==='error'||j.last_result==='fail';
h+='<div class="job'+(j.enabled?'':' disabled')+'"><div class="job-top"><div>';
h+='<div class="job-name">'+esc(j.name)+'</div>';
h+='<div class="job-schedule">'+esc(j.schedule||'manual trigger')+'</div>';
h+='</div><div class="job-actions">';
h+='<label class="toggle"><input type="checkbox" '+(j.enabled?'checked':'')+' onchange="tog(''+j.id+'')"><span class="sl"></span></label>';
h+='<button class="btn btn-sm btn-run" onclick="run(''+j.id+'')">Run Now</button>';
h+='<button class="btn btn-sm" onclick="openEdit(''+j.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+j.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';
if(j.command)h+='<div class="job-cmd">$ '+esc(j.command)+'</div>';
if(j.webhook_url)h+='<div class="job-cmd">&#8599; '+esc(j.webhook_url)+'</div>';
h+='<div class="job-meta">';
if(j.last_result)h+='<span class="'+(ok?'result-ok':fail?'result-fail':'')+'">Last: '+esc(j.last_result)+'</span>';
h+='<span>Runs: '+j.run_count+'</span>';
if(j.fail_count)h+='<span class="result-fail">Fails: '+j.fail_count+'</span>';
if(j.last_run_at)h+='<span>Last run: '+ft(j.last_run_at)+'</span>';
h+='</div></div>';
});
document.getElementById('list').innerHTML=h;
}

async function tog(id){await fetch(A+'/jobs/'+id+'/toggle',{method:'POST'});load();}
async function run(id){await fetch(A+'/jobs/'+id+'/run',{method:'POST'});load();}
async function del(id){if(!confirm('Delete this job?'))return;await fetch(A+'/jobs/'+id,{method:'DELETE'});load();}

function formHTML(job){
var i=job||{name:'',schedule:'',command:'',webhook_url:''};
var isEdit=!!job;
var h='<h2>'+(isEdit?'EDIT JOB':'NEW SCHEDULED JOB')+'</h2>';
h+='<div class="fr"><label>Name *</label><input id="f-name" value="'+esc(i.name)+'" placeholder="e.g. Database backup"></div>';
h+='<div class="fr"><label>Schedule (cron expression)</label><input id="f-schedule" value="'+esc(i.schedule)+'" placeholder="*/5 * * * *"><div class="hint">min hour dom mon dow &#8212; e.g. 0 2 * * * = daily at 2am</div></div>';
h+='<div class="fr"><label>Command (shell)</label><input id="f-cmd" value="'+esc(i.command)+'" placeholder="e.g. pg_dump mydb > /backups/db.sql"></div>';
h+='<div class="fr"><label>Webhook URL (alternative to command)</label><input id="f-webhook" value="'+esc(i.webhook_url)+'" placeholder="https://example.com/webhook"></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Create Job')+'</button></div>';
return h;
}

function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');document.getElementById('f-name').focus();}
function openEdit(id){var job=null;for(var j=0;j<jobs.length;j++){if(jobs[j].id===id){job=jobs[j];break;}}if(!job)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(job);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}

async function submit(){
var name=document.getElementById('f-name').value.trim();
if(!name){alert('Name is required');return;}
var body={name:name,schedule:document.getElementById('f-schedule').value.trim(),command:document.getElementById('f-cmd').value.trim(),webhook_url:document.getElementById('f-webhook').value.trim()};
if(editId){await fetch(A+'/jobs/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/jobs',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();
}

function ft(t){if(!t)return'';try{var d=new Date(t);return d.toLocaleDateString('en-US',{month:'short',day:'numeric'})+' '+d.toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
