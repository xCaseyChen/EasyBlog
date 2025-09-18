const formatter = new Intl.DateTimeFormat("en-CA", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit"
});

document.querySelectorAll('[datetime]').forEach(el => {
    const utcStr = el.getAttribute('datetime');
    const d = new Date(utcStr);
    el.textContent = el.classList.contains("created")
        ? "Created At: " + formatter.format(d)
        : "Updated At: " + formatter.format(d);
});