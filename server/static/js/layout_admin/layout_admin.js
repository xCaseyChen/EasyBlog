const logoutBtn = document.getElementById("logoutBtn");
const logoutDialog = document.getElementById("logoutDialog");
const logoutConfirmBtn = document.getElementById("logoutConfirmBtn")
const logoutCancelBtn = document.getElementById("logoutCancelBtn")

const resultDialog = document.getElementById("resultDialog");
const resultMessage = document.getElementById("resultMessage");
const resultOkBtn = document.getElementById("resultOkBtn");

function showResult(message, isSuccess) {
    resultMessage.textContent = message;

    resultDialog.showModal();

    resultOkBtn.onclick = () => {
        resultDialog.close();
        if (isSuccess) {
            window.location.href = "/home";
        }
    };
}

logoutBtn.addEventListener("click", async () => {
    logoutDialog.showModal();
});

logoutConfirmBtn.addEventListener("click", async () => {
    try {
        const resp = await fetch("/api/admin/logout", { method: "POST" });
        logoutDialog.close();
        if (resp.ok) {
            showResult("Logout success", true);
        } else {
            showResult("Logout failed", false);
        }
    } catch (err) {
        console.error(err);
        showResult("Network error", false);
    }
});

logoutCancelBtn.addEventListener("click", async () => {
    logoutDialog.close();
});