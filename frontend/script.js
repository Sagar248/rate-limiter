const API_URL = "http://localhost:8080/api"; // change after deploy

async function hitApi() {
  const userId = document.getElementById("userId").value || "demo";

  const res = await fetch(API_URL, {
    headers: { "X-User-ID": userId }
  });

  const data = await res.json();
  log(data);
}

async function burst() {
  for (let i = 0; i < 10; i++) {
    await hitApi();
  }
}

function log(data) {
  const output = document.getElementById("output");
  output.innerText = JSON.stringify(data, null, 2);
}