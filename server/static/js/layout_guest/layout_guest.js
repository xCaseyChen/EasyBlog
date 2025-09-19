const loginBtn = document.getElementById("loginBtn");
const loginDialog = document.getElementById("loginDialog");
const loginForm = document.getElementById("loginForm");
const loginSubmitBtn = loginForm.querySelector(".btn-submit");
const loginCancel = loginForm.querySelector(".btn-cancel");

const resultDialog = document.getElementById("resultDialog");
const resultMessage = document.getElementById("resultMessage");
const resultOkBtn = document.getElementById("resultOkBtn");

function showResult(message, isSuccess) {
    resultMessage.textContent = message;
    // resultMessage.style.color = isSuccess ? "green" : "red";

    resultDialog.showModal();

    resultOkBtn.onclick = () => {
        resultDialog.close();
        if (isSuccess) {
            window.location.href = "/admin/dashboard";
        }
    };
}

loginBtn.addEventListener("click", async () => {
    try {
        const resp = await fetch("/api/admin/ping", { method: "GET" });
        if (resp.ok) {
            window.location.href = "/admin/dashboard";
        } else if (resp.status === 401) {
            loginForm.password.value = "";
            loginDialog.showModal();
        } else {
            showResult("Login failed", false);
        }
    } catch (err) {
        console.error(err);
        showResult("Network error", false);
    }
});

loginCancel.addEventListener("click", () => {
    loginForm.password.value = "";
    loginDialog.close();
});

loginSubmitBtn.addEventListener("click", async (e) => {

    e.preventDefault();

    loginSubmitBtn.disabled = true;

    const password = loginForm.password.value.trim();
    if (!password) {
        showResult("Please enter password", false);
        loginSubmitBtn.disabled = false;
        return;
    }

    try {
        const response = await fetch("/api/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ password }),
        });

        const data = await response.json();

        if (response.ok && data.success) {
            loginDialog.close();
            showResult("Login success", true);
        } else {
            showResult(data.message || "Login failed", false);
        }
    } catch (err) {
        console.error(err);
        showResult("Network error", false);
    } finally {
        loginSubmitBtn.disabled = false;
    }
});