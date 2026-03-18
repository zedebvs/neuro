# 📦 Структура проекта (Project Layout)

### 📂 Основные пакеты
<table>
  <thead>
    <tr>
      <th>Путь</th>
      <th>Назначение</th>
      <th>Ответственный пакет</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>/app</code></td>
      <td>Точки входа в приложение (main.go)</td>
      <td><code>package main</code></td>
    </tr>
    <tr>
      <td><code>/core</code></td>
      <td>Ядро для оркестра над апишками - крч маршрутизатор и регулятор</td>
      <td><code>package service</code>, <code>repository</code></td>
    </tr>
    <tr>
      <td><code>/audio</code></td>
      <td>Все что связано со звуком</td>
      <td><code>package</code>, <code></code></td>
    </tr>
    <tr>
      <td><code>/lin</code></td>
      <td>Абстракция над нейронками + парсер ответов и отправка в сервис, мб потом по ддругому сделаю</td>
      <td><code>package</code>, <code></code></td>
    </tr>
    <tr>
      <td><code>/network</code></td>
      <td>сеть + ядро encTCP для взаимодействия с микросервисом</td>
      <td><code>package</code>, <code></code></td>
    </tr>
    <tr>
      <td><code>/memory</code></td>
      <td>контроллер памятии - пока 4 слоя, разделенных на 2</td>
      <td><code>package</code>, <code></code></td>
    </tr>
  </tbody>
</table>

---

### 📝 Подробное описание компонентов

<details>
<summary><b>🔍 Нажми, чтобы развернуть описание API</b></summary>

Здесь можно описать роуты или специфику работы `internal/api`:
- Используется фреймворк **Gin/Chi**.
- Авторизация через **JWT**.

</details>

<details>
<summary><b>🗄️ Нажми, чтобы развернуть описание DB</b></summary>

- Слой миграций лежит в `/migrations`.
- Для работы с БД используется **GORM** или **sqlx**.

</details>

---