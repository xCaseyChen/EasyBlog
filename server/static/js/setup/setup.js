const msgDialog = document.getElementById("msgDialog");
const msgContent = document.getElementById("msgContent");
const msgOk = document.getElementById("msgOk");

function showMessage(message, callback) {
    msgContent.textContent = message;
    msgDialog.showModal();

    function handleOk() {
        msgDialog.close();
        msgOk.removeEventListener("click", handleOk);
        if (callback) callback();
    }

    msgOk.addEventListener("click", handleOk);
}

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
            showMessage(data.message, () => {
                window.location.href = "/home";
            });
        } else {
            showMessage(data.message || "Request failed");
        }
    } catch (err) {
        showMessage("Network error");
        console.error(err);
    }
});