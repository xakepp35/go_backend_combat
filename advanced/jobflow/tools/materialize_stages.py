#!/usr/bin/env python3
"""Rebuild open milestone snapshots from the canonical reference and explicit early implementations."""
import json
from pathlib import Path
import re
import shutil
import subprocess
import os

ROOT=Path(__file__).resolve().parents[1]
FINAL=ROOT/'reference/24'

def selected(source,names):
    matches=list(re.finditer(r'^func\s+(?:\([^)]*\)\s*)?(\w+)\s*\(',source,re.M))
    if not matches:return source
    out=source[:matches[0].start()]
    for i,m in enumerate(matches):
        if m.group(1) in names:
            out+=source[m.start():matches[i+1].start() if i+1<len(matches) else len(source)]
    return imports(out)

def imports(source):
    block=re.search(r'(?m)^import \(\n(.*?)\n\)',source,re.S)
    if not block:return source
    body=source[:block.start()]+source[block.end():]
    tick=chr(96)
    body=re.sub(r'(?s)/\*.*?\*/|//[^\n]*|"(?:\\.|[^"\\])*"|'+tick+'[^'+tick+']*'+tick,'',body)
    keep=[]
    for line in block.group(1).splitlines():
        m=re.search(r'"([^"]+)"',line)
        if not m:continue
        path=m.group(1)
        pkg={'github.com/jackc/pgx/v5':'pgx','github.com/redis/go-redis/v9':'redis','github.com/segmentio/kafka-go':'kafka'}.get(path,path.split('/')[-1])
        if re.search(r'\b'+re.escape(pkg)+r'\.',body):keep.append(line)
    return source[:block.start()]+('import (\n'+'\n'.join(keep)+'\n)' if keep else '')+source[block.end():]

def write(root,path,source):
    p=root/path;p.parent.mkdir(parents=True,exist_ok=True);p.write_text(source)

def copy(root,path):
    p=root/path;p.parent.mkdir(parents=True,exist_ok=True);shutil.copy2(FINAL/path,p)

def strip_handler(source,pattern):
    start=source.index('\tm.HandleFunc("'+pattern+'"')
    end=source.index('\n\t})',start)+len('\n\t})')
    return source[:start]+source[end:]

def main():
    go=os.environ.get('GO','go')
    lessons=json.loads((ROOT/'course/lessons.json').read_text())
    base=FINAL/'internal/task/task.go'
    # Early task model intentionally excludes UUID/framing/lease concepts.
    simple='''package task
import("errors";"time")
type State string
const Draft State="DRAFT"
var ErrInvalid=errors.New("invalid task")
var ErrNotFound=errors.New("task not found")
var ErrConflict=errors.New("task conflict")
var ErrNotReady=errors.New("task not ready")
var ErrExpired=errors.New("task expired")
var ErrStale=errors.New("stale attempt")
var ErrQuota=errors.New("quota exceeded")
type Task struct {
 ID string `json:"id"`; Kind string `json:"kind"`; Body []byte `json:"body"`;
 State State `json:"state"`; Version int64 `json:"version"`;
 NBF uint32 `json:"nbf,omitempty"`; EXP uint32 `json:"exp,omitempty"`;
 UpdatedAt time.Time `json:"updated_at"`
}
func New(id,kind string,body []byte)Task{return Task{ID:id,Kind:kind,Body:append([]byte(nil),body...),State:Draft,Version:1}}
func Clone(t Task)Task{t.Body=append([]byte(nil),t.Body...);return t}
'''
    validate='''func(t Task)Validate()error{if t.ID==""||t.Kind==""||len(t.Body)>1<<20{return ErrInvalid};if t.EXP!=0&&t.NBF>=t.EXP{return ErrInvalid};return nil}
'''
    for n in range(1,24):
        root=ROOT/'reference'/f'{n:02}'
        if root.exists():shutil.rmtree(root)
        root.mkdir()
        write(root,'go.mod','module jobflow\n\ngo 1.27.0\n')
        for p in ['cmd/hello/main.go','.github/workflows/ci.yml','Makefile']:
            if (FINAL/p).exists():copy(root,p)
        if n<11:
            p=root/'.github/workflows/ci.yml';p.write_text(p.read_text().replace('go test -race','go test'))
            p=root/'Makefile';p.write_text(p.read_text().replace('test -race','test'))
        write(root,'.lesson',f'{n:02}\n')
        if n>=2:
            write(root,'internal/task/task.go',simple+(validate if n>=4 else 'func(t Task)Validate()error{return nil}\n'))
            write(root,'cmd/pub/main.go','''package main
import("fmt";"jobflow/internal/task")
func main(){a:=task.New("a","echo",[]byte("report"));fmt.Println(a.ID,a.Kind,string(a.Body))}
''')
        if n>=3:
            source=(FINAL/'internal/memory/repository.go').read_text()
            if n<6:
                source=re.sub(r'\n\s*mu\s+sync.Mutex','',source)
                source=re.sub(r'\n\s*r.mu.Lock\(\)','',source)
                source=re.sub(r'\n\s*defer r.mu.Unlock\(\)','',source)
            write(root,'internal/memory/repository.go',imports(source))
        if n>=5:copy(root,'internal/app/service.go')
        if n>=6:
            copy(root,'internal/api/api.go')
            write(root,'cmd/pub/main.go','''package main
import("log";"net/http";"os";"time";"jobflow/internal/api";"jobflow/internal/memory")
func main(){addr:=os.Getenv("HTTP_ADDR");if addr==""{addr=":8080"};s:=&http.Server{Addr:addr,Handler:api.New(memory.New()),ReadHeaderTimeout:5*time.Second};log.Fatal(s.ListenAndServe())}
''')
        if n>=7:
            copy(root,'migrations/001_core.sql')
            write(root,'sql/01_basics.sql',"INSERT INTO drafts(id,doc) VALUES('example','{\"kind\":\"echo\"}');\nSELECT id,doc FROM drafts ORDER BY id;\nUPDATE drafts SET doc='{\"kind\":\"sum\"}' WHERE id='example';\nDELETE FROM drafts WHERE id='example';\n")
        if n>=8:
            copy(root,'internal/postgres/store.go')
            write(root,'cmd/pub/main.go','''package main
import("context";"log";"net/http";"os";"time";"jobflow/internal/api";"jobflow/internal/postgres")
func main(){s,err:=postgres.Open(context.Background(),os.Getenv("DATABASE_URL"));if err!=nil{log.Fatal(err)};defer s.Close();if err=s.Migrate(context.Background(),"migrations");err!=nil{log.Fatal(err)};addr:=os.Getenv("HTTP_ADDR");if addr==""{addr=":8080"};server:=&http.Server{Addr:addr,Handler:api.New(s),ReadHeaderTimeout:5*time.Second};log.Fatal(server.ListenAndServe())}
''')
        if n>=9:copy(root,'internal/postgres/intent.go')
        if n>=10:
            copy(root,'internal/pool/pool.go')
            if n==10:
                write(root,'internal/pool/pool.go','''package pool
import("context";"errors";"sync")
func Run(ctx context.Context,workers int,jobs <-chan func(context.Context))error{if workers<1||workers>256{return errors.New("workers must be 1..256")};var wg sync.WaitGroup;for range workers{wg.Add(1);go func(){defer wg.Done();for f:=range jobs{f(ctx)}}()};wg.Wait();return nil}
''')
        if n>=11:
            copy(root,'internal/api/server.go')
            write(root,'cmd/pub/main.go','''package main
import("context";"log";"os";"os/signal";"syscall";"jobflow/internal/api";"jobflow/internal/postgres")
func main(){ctx,cancel:=signal.NotifyContext(context.Background(),os.Interrupt,syscall.SIGTERM);defer cancel();s,err:=postgres.Open(ctx,os.Getenv("DATABASE_URL"));if err!=nil{log.Fatal(err)};defer s.Close();if err=s.Migrate(ctx,"migrations");err!=nil{log.Fatal(err)};addr:=os.Getenv("HTTP_ADDR");if addr==""{addr=":8080"};if err=api.RunHTTP(ctx,addr,api.New(s));err!=nil{log.Fatal(err)}}
''')
        if n>=12:copy(root,'internal/protocol/header.go')
        if n>=13:
            copy(root,'internal/protocol/frame.go');copy(root,'internal/protocol/frame_test.go')
        if n>=14:
            copy(root,'internal/task/task.go');copy(root,'internal/task/graph.go');copy(root,'migrations/002_engine.sql')
            funcs={'stripped','digest','putPayload','saveJob','loadJob','Job','Jobs','Payload','Submit'}
            if n>=19:funcs|={'depsDone','Claim','Complete','Effect'}
            if n>=20:funcs.add('Maintain')
            write(root,'internal/postgres/engine.go',selected((FINAL/'internal/postgres/engine.go').read_text(),funcs))
            names={'Events','Delivered'}|({'GC','Collect'} if n>=22 else set())
            write(root,'internal/postgres/events.go',selected((FINAL/'internal/postgres/events.go').read_text(),names))
        if n>=15:copy(root,'internal/cache/cache.go')
        if n>=16:
            names={'PublishOnce','Relay'}|({'ConsumeOne','Consume'} if n>=17 else set())
            write(root,'internal/bus/kafka.go',selected((FINAL/'internal/bus/kafka.go').read_text(),names))
        if n>=17:
            names={'Key','Apply','Get'}|({'Ready','Group','Ack','Receive'} if n>=18 else set())
            write(root,'internal/projection/projection.go',selected((FINAL/'internal/projection/projection.go').read_text(),names))
            copy(root,'internal/client/client.go');copy(root,'internal/process/process.go')
            if n==17:
                source=(FINAL/'cmd/sub/main.go').read_text()
                start=source.index('\t\t\t\tif e = p.Ready(')
                end=source.index('\n\t\t\t\t}',start)+len('\n\t\t\t\t}')
                source=source[:start]+source[end:]
                # Retain the read to detect a missing hash, before Streams scheduling exists.
                source=source.replace('t, e := p.Get(ctx, id)','_, e := p.Get(ctx, id)')
                write(root,'cmd/sub/main.go',imports(source))
        if n>=18:copy(root,'cmd/sub/main.go')
        if n>=19:
            copy(root,'internal/runner/runner.go');copy(root,'cmd/act/main.go');copy(root,'migrations/003_receipts.sql')
        if n>=14:
            source=(FINAL/'internal/api/engine.go').read_text()
            if n<22:source=strip_handler(source,'POST /internal/gc')
            if n<23:source=strip_handler(source,'GET /metrics')
            if n<19:
                for route in ['POST /tasks/{id}/claim','POST /tasks/{id}/result','POST /tasks/{id}/effects']:
                    source=strip_handler(source,route)
            if n==14:
                source=strip_handler(source,'GET /payloads/{key}')
                source=source.replace(', cas cache.CAS, retention time.Duration','')
            write(root,'internal/api/engine.go',imports(source))
            copy(root,'internal/process/process.go')
            source=(FINAL/'cmd/pub/main.go').read_text()
            if n<20:
                start=source.index('\tgo func() {\n\t\tdefer wg.Done()\n\t\tticker')
                end=source.index('\n\terr = api.RunHTTP',start)
                source=source[:start]+source[end:]
                source=source.replace('wg.Add(2)','wg.Add(1)')
            if n<16:
                start=source.index('\tvar wg sync.WaitGroup')
                end=source.index('\n\terr = api.RunHTTP',start)
                source=source[:start]+source[end:]
                source=source.replace('\twg.Wait()\n','')
            if n==14:
                start=source.index('\tr := process.Redis()')
                end=source.index('\tm := api.New(s)',start)
                source=source[:start]+source[end:]
                source=source.replace('api.RegisterEngine(m, s, cas, time.Hour)','api.RegisterEngine(m, s)')
            write(root,'cmd/pub/main.go',imports(source))
        if n>=22:copy(root,'cmd/ret/main.go')
        if n>=23:
            copy(root,'Dockerfile');copy(root,'cmd/demo/main.go')
        # The module versions are locked, including transitive checksums. No @latest in student workflow.
        if n>=8:
            mod=(FINAL/'go.mod').read_text()
            write(root,'go.mod',mod)
            copy(root,'go.sum')
        run=subprocess.run([go,'mod','tidy'],cwd=root,capture_output=True,text=True)
        if run.returncode:raise RuntimeError(f'stage {n:02}: {run.stderr[-2000:]}')
        subprocess.run([str(Path(go).with_name('gofmt')),'-w',str(root)],check=True)
        row=lessons[n-1]
        write(root,'README.md',f'# JobFlow — этап {n:02}: {row["title"]}\n\nЭто открытый эталон **после** выполнения ДЗ {n:02}. Сначала решите работу самостоятельно.\n\nИнвариант: {row["invariant"]}\n\nИз корня курса: make check LESSON={n:02} PROJECT=reference/{n:02}\n\nРанние cmd/pub демонстрируют текущую границу; приём workflow по HTTP появляется на этапе 14, relay — на этапе 16, исполнение — на этапе 19, Compose — на этапе 23. Наличие полей будущих состояний с этапа 14 сохраняет общий контракт, но не делает будущие сервисы обязательными заранее.\n')
        print(f'{n:02} materialized',flush=True)
    write(FINAL,'.lesson','24\n')

if __name__=='__main__':main()
