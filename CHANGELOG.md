### v5.2.1
* обновлены зависимости
### v5.2.0
* Обновлены зависимости
* Добавлены новые endpoint'ы: `/application/update_application`, `/application/create_application`, `/application_group/create`, `/application_group/update`, `/application_group/delete_list`, `/application_group/get_by_id_list`, `/application_group/get_all`
### v5.1.1
* обновлены зависимости
### v5.1.0
* все сервисы (группы приложений) перенесены в домен с id = 1
* при создании/изменении сервисов (групп приложений) будет использоваться только домен с id = 1

### v5.0.0
* изменен механизм инициализации токена панели администратора
  * для механизма инициализации используется поле `baseline.initialAdminUiToken` в конфигурации
    * если поле не заполнено, механизм не запускается
    * иначе проверяется существование токена в базе, если токен уже создан, механизм не запускается
    * иначе создаются все сущности и выпускается бессрочный токен, приложению выдается доступ к методу `admin/auth/login`
    * механизм запускается при получении конфигурации
    * при конкурентном запуске модуля, механизм запускается только на одном инстансе
    * логи механизма маркируются полем  "worker": "baseline"
* обновлена зависимости
### v4.0.0
* Удален `InstanceUuid`
* Удалена интеграция с Redis
* Генерация токенов теперь использует обычный crypto/rand вместо JWT
* Убран из конфига defaultTokenExpireTime
* Добавлены методы проверки подлинности токенов и прав у приложений
  * system/secure/authenticate
  * system/secure/authorize
### v3.0.2
* исправлена ошибка при обновлении имени или описании у созданного application
### v3.0.1
* mark `InstanceUuid` as required
### v3.0.0
* migrate to isp-kit
* split domain and entity structs
* prepare to remove system block
  * remove CRUD for system
  * use default system_id
* remove migration initializing city_module
* unite requests to redis into one directory
  * prepare to remove redis
* update go to 1.17
### v2.2.7
* updated dependencies
* migrated to common local config
### v2.2.6
* fix migrations
### v2.2.5
* updated dependencies
### v2.2.4
* updated isp-lib
### v2.2.3
* updated isp-lib
* updated isp-lib-test
### v2.2.2
* updated isp-lib
* updated isp-event-lib
### v2.2.1
* fix linter
### v2.2.0
* update libs
### v2.1.0
* update to go mod
### v2.0.0
* update `isp-lib` to 2.0.0
### v1.2.0
* add `access list` methods
* fix revoke token
### v1.1.3
* update to new log
* migrate to new db client
### v1.1.2
* add document generation
### v1.1.1
* update config description
* update lib
### v1.1.0
* add default remote configuration
