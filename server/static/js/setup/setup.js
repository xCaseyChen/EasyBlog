document.getElementById("submitBtn").addEventListener("click", async () => {
    const password = document.getElementById("password").value;

    try {
        const response = await fetch("/api/setup/password", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ password }),
        });

        const data = await response.json();
        if (response.ok && data.success) {
            alert(data.message);
            window.location.href = "/home";
        } else {
            alert(data.message || "Request failed");
        }
    } catch (err) {
        alert("Network error");
        console.error(err);
    }
});