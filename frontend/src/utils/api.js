export async function postJSON(url, body) {
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  if (!res.ok) {
    let err;
    try {
      err = await res.json();
    } catch {
    }
    throw new Error(err?.error || `Request failed: ${res.status}`);
  }

  try {
    return await res.json();
  } catch {
    return {};
  }
}
