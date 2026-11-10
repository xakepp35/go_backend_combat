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

Задачи и результаты — только в корневой `.zts/{pub,sub,map,act}` проекта.
Использовать один скилл/бинарник `zts` и ровно четыре макрокоманды: pub — контракт
и порядок, sub — состояние/контекст, map — подготовка, act — проверки/фиксация.
Перед работой читать `zts sub` и `zts sub --next`. Сначала декомпозиция; активная
задача всегда первая. Progress хранить через `zts sub --progress`, результаты —
в act. Очередь head/tail принадлежит sub. AGENTS.md содержит устойчивые правила,
без текущего статуса и истории. Отдельные checkpoint/changelog/continuation и
обязательные промпты передачи не создавать. Очередь не расширяет разрешений,
не запускает goal, субагентов или фонового исполнителя. Один исполнитель на
проект; чужие изменения/staging сохранять. Не выдавать исторические проверки
за новые и не подменять машинную приёмку словами агента.
