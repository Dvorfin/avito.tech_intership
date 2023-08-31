# Приложение-интерфейс для работы с Redis'ом с проксированием трафика.

Использованные технологии:
- redis (в качестве хранилища данных)
- golang (для создания приложения взаимдействия)
- nginx (для проксирования трафика к приложению)
- Docker (для контейнеризации компонентов и их запуска)

# Настройки
При решении задания возникли сомнения в формулировке "приложение-интерфейс". Оно должно представлять из себя просто возможность отправлять, например, curl запросы на эндпоинты или должно иметь некоторый ui? По этой причине реализовано две версии: с поддержкой веб-интерфейса и без. Также было непонятно, что подразумевается под "конфигурация не запекается внутри образов". Решил, что имеется ввиду наличие параметров ENV в докер файлах, для возоможности их изменения. Поэтому в dockerfile'ы redis'a и приложения добавлены поля `REDIS_PASSWORD`, `REDIS_TLS_PORT` и `WITH_UI` соответственно.

По дефолту приложение запсукается с веб-интерфейсом. Для того чтобы запустить его без - в файле конфигурации `envfile.env` нужно удалить или изменить значение `WITH_UI`.

Также в файле конфигурации `envfile.env` можно задать пароль для redis'a в параметре `REDIS_PASSWORD`. При изменении пароля redis'а требуется установить заданный пароль в приложении в файле config.yaml (app/src/config/config.yaml) параметр `pass`. Помимо этого в `config.yaml` указан параметр `address`, отвечающий за ip адрес контейнера с запущенным redis. 

На redis'е добавлена аутентификация по паролю и взаимодействие с приложением через TLS. Для этого созданы самоподписанные SSL сертефикаты ([link](https://www.youtube.com/watch?v=SZENu4kzZe0)). Сертефикаты приложения лежат по пути` app/src/tests
/tls/`, сертефикаты reedis'a в `redis/tls/` соответственно.

## Запуск с веб-интерфейсом <a name="with_ui"></a>
```sh
docker compose up
```
После успешного запуска перейти по адресу tech-avito-test:8089, откроется окно:

![image](https://github.com/Dvorfin/avito.tech_intership/assets/70969469/619c5c82-1ee2-45cf-8974-056c950f2d06)


> Если приложение не доступно, то на хостовой машине требуется добавить в файл /etc/hosts следующие строки:
> 127.0.0.1 avito-tech
> 127.0.0.1 tech-avito-test
tech-avito-test и avito-tech - соответственно адреса по которым nginx проксирует трафик в приложение


Веб-интерфейс реализует:

1) Получение значения по ключу.
- в случае наличия ключа в БД будет получен респонс:
```json
{
  "Key": "введенный ключ"
  "Value": "значение ключа"
}
```
- в случае отсутствия ключа:
```json
{
  "Msg": "There is no such key."
}
```

2) Установку ключа.
- в случае отсутствия ключа в БД будет получен респонс:
```json
{
  "Msg": "Key setted."
}
```

- в случае попытки установить/обновить ключ пустым значением будет получен респонс:
```json
{
  "Msg": "Key updated. Notice, that value was empty!"
}
```

3) Удаление ключа.
- в случае наличия ключа в БД будет получен респонс:
```json
{
  "Msg": "Key sucessfully deleted!"
}
```

- в случае отсутствия ключа:
```json
{
  "Msg": "There is no such key."
}
```


## Запуск без веб-интерфейса
```sh
docker compose up
```
Для взаимодействия воспользуемся curl.
1) Для получения значения по ключу:
```curl
curl -i tech-avito-test:8089/get_key?key=<указать значение ключа>
```

2) Для установки/обновления ключа:
```curl
curl -X POST tech-avito-test:8089/set_key -d '{"key": "<ключ>", "value": "<значение>"}'
```
Важно! Ключ и значение должны быть строкой!
`<ключ>` заменить на ключ без `<` `>`, например: `<ключ>` -> `myKey`
 
- В случае попытки отправить другой тип ключа будет получено сообщение:
```json
{
  "Msg": "Invalid key type! Must be string."
}
```
- В случае некорректного значения:
```json
{
  "Msg": "Invalid value type! Must be string."
}
```

3) Для удаления ключа:
```curl
curl -i -X DELETE tech-avito-test:8089/del_key -d '{"<ключ>": "<значение>"}'
```

Все респонсы аналогичны описанным в блоке [Запуск с веб-интерфейсом](#with_ui).


## TODO
Вещи, которые можно улучшить/доработать:
1) Добавить пребилд контейнер при сборке приложения на golang для уменьшения веса конечного контейнера.
2) Вынести повторяющийся код в приложении в отдельные модули (например подключение к БД и проверку SSL ключей).
2) Доработать конфиги как в приложении, так и в docker compose (добавить настройку ip и портов, поработать над более оптимальным заданием env)
3) 

