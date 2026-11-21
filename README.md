# ระบบ Login - Register

โปรเจกต์ตัวอย่างสมัครสมาชิกและเข้าสู่ระบบด้วย Vue, Go, JWT และ PostgreSQL มีไว้สำหรับทดลองและเรียนรู้การเชื่อมต่อ frontend กับ backend

## Tech Stack

- **Frontend:** Vue 3, Vue Router, TypeScript, Vite, Tailwind CSS, Axios
- **Backend:** Go, Fiber, pgx, JWT, godotenv
- **Database:** PostgreSQL
- **เครื่องมือ:** Docker Compose, Playwright, Prettier, vue-tsc

## ก่อนเริ่ม

ต้องมี Go 1.26.2+, Node.js 22.18+ หรือ 24.12+, npm และ Docker Compose

คัดลอกไฟล์ตั้งค่าทั้งสองชุด:

```bash
cp Backend/.env.example Backend/.env
cp Frontend/.env.example Frontend/.env
```

ตั้งค่า API URL ใน `Frontend/.env` เป็น `API_URL=http://localhost:8080/api/v1` และกำหนด `JWT_SECRET` กับรหัสผ่านฐานข้อมูลใน `Backend/.env` ค่าตัวอย่างเหมาะสำหรับเครื่อง local เท่านั้น

## วิธีรัน

เปิดฐานข้อมูล:

```bash
docker compose -f Backend/docker-compose.yml up -d db
```

เปิด backend ใน terminal หนึ่ง:

```bash
cd Backend
go run ./cmd/api
```

เปิด frontend อีก terminal:

```bash
cd Frontend
npm install
npm run dev
```

เปิดเว็บที่ <http://localhost:5173> และตรวจ API ได้ที่ <http://localhost:8080/api/v1/health>
