package server
import "net/http"
func(s *Server)dashboard(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html");w.Write([]byte(dashHTML))}
const dashHTML=`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Mainspring</title>
<style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}
.main{padding:1.5rem;max-width:900px;margin:0 auto}
.job{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem}
.job-top{display:flex;justify-content:space-between;align-items:center}
.job-name{font-size:.82rem;color:var(--cream)}
.job-schedule{font-size:.65rem;color:var(--gold);margin-top:.1rem}
.job-cmd{font-size:.65rem;color:var(--cd);background:var(--bg);padding:.2rem .4rem;border:1px solid var(--bg3);margin-top:.3rem;word-break:break-all}
.job-meta{font-size:.6rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.8rem}
.job-result{font-size:.6rem;margin-top:.2rem;padding:.2rem .4rem}
.result-success{background:#4a9e5c22;color:var(--green);border:1px solid #4a9e5c44}
.result-failed{background:#c9444422;color:var(--red);border:1px solid #c9444444}
.toggle{position:relative;width:36px;height:18px;cursor:pointer;display:inline-block;vertical-align:middle}
.toggle input{opacity:0;width:0;height:0}.toggle .sl{position:absolute;inset:0;background:var(--bg3);border-radius:9px;transition:.2s}
.toggle .sl:before{content:'';position:absolute;width:14px;height:14px;left:2px;bottom:2px;background:var(--cm);border-radius:50%;transition:.2s}
.toggle input:checked+.sl{background:var(--green)}.toggle input:checked+.sl:before{transform:translateX(18px);background:var(--cream)}
.btn{font-size:.6rem;padding:.25rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:420px;max-width:90vw}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.5rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.15rem}
.fr input,.fr select{width:100%;padding:.35rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1>MAINSPRING</h1><button class="btn btn-p" onclick="openForm()">+ New Job</button></div>
<div class="main" id="main"></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let jobs=[];
async function load(){const r=await fetch(A+'/jobs').then(r=>r.json());jobs=r.jobs||[];render();}
function render(){if(!jobs.length){document.getElementById('main').innerHTML='<div class="empty">No scheduled jobs. Create one to get started.</div>';return;}
let h='';jobs.forEach(j=>{
h+='<div class="job"><div class="job-top"><div><div class="job-name">'+esc(j.name)+'</div><div class="job-schedule">'+esc(j.schedule||'manual')+'</div></div><div style="display:flex;gap:.4rem;align-items:center"><label class="toggle"><input type="checkbox" '+(j.enabled?'checked':'')+' onchange="tog(\''+j.id+'\')"><span class="sl"></span></label><button class="btn" onclick="run(\''+j.id+'\')">Run Now</button><button class="btn" onclick="del(\''+j.id+'\')" style="color:var(--cm)">✕</button></div></div>';
if(j.command)h+='<div class="job-cmd">$ '+esc(j.command)+'</div>';
if(j.webhook_url)h+='<div class="job-cmd">→ '+esc(j.webhook_url)+'</div>';
h+='<div class="job-meta"><span>Runs: '+j.run_count+'</span><span>Fails: '+j.fail_count+'</span>';
if(j.last_run_at)h+='<span>Last: '+ft(j.last_run_at)+'</span>';h+='</div>';
if(j.last_result)h+='<div class="job-result '+(j.last_result.startsWith('error')||j.last_result.startsWith('fail')?'result-failed':'result-success')+'">'+esc(j.last_result)+'</div>';
h+='</div>';});
document.getElementById('main').innerHTML=h;}
async function tog(id){await fetch(A+'/jobs/'+id+'/toggle',{method:'POST'});load();}
async function run(id){await fetch(A+'/jobs/'+id+'/run',{method:'POST'});load();}
async function del(id){if(confirm('Delete?')){await fetch(A+'/jobs/'+id,{method:'DELETE'});load();}}
function openForm(){document.getElementById('mdl').innerHTML='<h2>New Scheduled Job</h2><div class="fr"><label>Name</label><input id="f-n" placeholder="e.g. Database backup"></div><div class="fr"><label>Schedule (cron)</label><input id="f-s" placeholder="*/5 * * * * (every 5 min)"></div><div class="fr"><label>Command (shell)</label><input id="f-c" placeholder="e.g. pg_dump mydb > backup.sql"></div><div class="fr"><label>Webhook URL (alternative)</label><input id="f-w" placeholder="https://..."></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="sub()">Create</button></div>';document.getElementById('mbg').classList.add('open');}
async function sub(){await fetch(A+'/jobs',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('f-n').value,schedule:document.getElementById('f-s').value,command:document.getElementById('f-c').value,webhook_url:document.getElementById('f-w').value})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function ft(t){if(!t)return'';return new Date(t).toLocaleDateString()+' '+new Date(t).toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'});}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();
</script></body></html>`
