# jog implementation invariants

- Go 1.27.3; standard library only. Public handler: func(context.Context, Q) (R, error).
- Bound goroutines and admission. Never close a queue while a sender can access it.
- Cancellation is cooperative; wait for admitted work before returning from a workflow.
- ContinueIndependent is the default: terminal node errors block only descendants.
  FailFast is explicit. Journal/codec/retry-callback/admission failures and external
  Run context cancellation stop the entire execution; other Runs remain isolated.
- Persist node success before releasing dependent nodes. Journal failures stop scheduling.
- Replay is at least once for unfinished nodes; never claim exactly-once external effects.
- Workflow definitions are immutable after Build. Read only declared dependencies.
- Keep examples executable. Verify concurrency with go test -race and testing/synctest.

## Управление работой

Задачи, незавершённый progress и результаты хранятся только в корневой
`.zts/{plan,done,head,tail}` этого Git-репозитория. Использовать общие `zplan` и
`zdone`; перед продолжением читать очередь и ближайшую задачу. В AGENTS.md —
только устойчивые правила, без текущего статуса и истории. Отдельные checkpoint,
changelog, continuation envelope и обязательные промпты передачи не создавать.
Очередь не расширяет разрешений и не запускает goal, субагентов или фонового
исполнителя. Работает один исполнитель; посторонние изменения/staging сохранять.
