export async function GET() {
  try {
    const res = await fetch('http://localhost:8080/ping');
    return new Response(res.body, { 
      status: res.status, 
      headers: res.headers 
    });
  } catch (err) {
    return new Response('Ошибка подключения к бэкенду', { status: 500 });
  }
}