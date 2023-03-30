## Описание проекта
Web-crawler - серверное приложение, которое обходит web-сайты, сканируя ссылки на них. Он начинает с какого-то корневого адреса (который задается вручную), а далее переходит по обнаруженным ссылкам (аналогично BFS на графе).

Для реализации распределенного варианта, будет использоваться Kafka. Для старта, в приложение будет подан изначальный URL, который попадет в topic X Кафки. Один из инстансов приложения прочитает его, скачает web-документ и просканирует все ссылки в нем. Все найденные ссылки он положит в тот же topic X, а сам web-документ в topic Y с скачанными документами.
Результирующий Kafka topic с документами с помощью File System Kafka Connect попадет, для простоты, в локальную файловую систему (есть возможность перелить их также в Hadoop, S3 и пр. распред. системы).

## Использованные технологии
1. Go 1.19
2. Kafka
3. Kafka Connect
4. Apache Camel Kafka Connectors
5. Postgresql

## Как это запустить?
1. Установить на машину Docker
2. Установить Docker-compose
3. Выполнить ``docker-compose up -d`` в корневой папке с проектом
4. Для добавления URL в очередь crawl нужно выполнить POST запрос на адрес http://localhost:8080/urls, в Body отправить URL (просто текстом).
5. В результате: crawler начнет обходить сайт, добавлять скачанные документы в папку .docker/docs в виде {domain}-{UUID}.html.