# Instructions for initial setup
  ```
- до поездки
	- узанать: android или ios
	- настроить точку в CRM
		- ТИС и.д

- до закрытия точки (за пол часа уже в заведении)
	- проверить инет
		- SpeedTest
	- почта
		- добавляем в тестеры в play market
	- если принтер по wifi
		- сфотать настройки принтера

- после закрытия точки		
	- устанавливаем приложение (по ссылки)
	- логинимся и настраиваем принтер
		- провести тестовую продажу
	- краткий инструктаж 
		- видосики на qr
  ```

# Coins logic
```
Логика Zebra Coin в мобильном приложении:

При регистрации 10 коинов пользователю
При каждой покупке с суммой свыше 2500 тенге, 10 коинов пользователю 
При каждом оставлении коментариев 5 коинов

Награды: 
1000 коинов - журнал 
1500 коинов - термокружка
2000 коинов - мерч 
2500 коинов - кофта
```


# ZebraCrm links


| Api | Url  Staging | Url  Production|
| ------ | ------ | ------ |
|   fronend    |   https://zebra.korsetu.kz  | https://zebra-crm.kz   |
|   backend   |    https://zebra-api.korsetu.kz/item/getAll    | https://zebra-crm.kz:8029/item/getAll    |
|   mobile-api   |    https://zebra-mobile-api.korsetu.kz/swagger    | https://zebra-crm.kz:13930/swagger    |
|   excel-gen-api   |    https://zebra-excel-gen-api.korsetu.kz/swagger    | https://zebra-crm.kz:31856/swagger    |


# Mobile Apps links

| App | Android | IOS | onelink |
| ------ | ------ | ------ | ------ |
|   ZebraCoffee    |   https://play.google.com/store/apps/details?id=kz.zebracrm.mobile_web  | https://apps.apple.com/ru/app/zebra-coffee/id6448264439   | https://onelink.to/xnfybg |
|   ZebraCoffee Terminal   |   https://play.google.com/store/apps/details?id=kz.zebracrm.mobile_terminal    | -    | - |

![qr__2_](/uploads/86d3e42b6201aa747f4e235427f759b4/qr__2_.png) **<-- ссылка на мобилку**



# A performance dashboard for Postgres - PGHero
 - https://zebra-crm.kz:15129/queries
 - user: `pg_hero_web_user`
 - pwd: see ci/cd vars - `PGHeroPwd`


# App Links and Universal Links urls

|  Android| IOS|
| ------ | ------ |
| https://zebra.korsetu.kz/.well-known/assetlinks.json  | https://zebra.korsetu.kz/.well-known/apple-app-site-association |
| https://zebra-crm.kz/.well-known/assetlinks.json    | https://zebra-crm.kz/.well-known/apple-app-site-association    |


<hr>
Важно при деплое на прод:

- Проверить работоспособность на **тесте** или **локально**
  - Если это новая формачка, то протестировать дабл клики сабмита 
- Если изменения касаются **UX/UI** то необходимо протестировать их на **7, 8, 10, 12** дюймовых планшетах (в браузере можно выставлять их) 
- При добавлении списка или таблицы для клиентов сервиса, добавить поиск и фильтрацию

