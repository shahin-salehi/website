function toggleCV() {
    const section = document.getElementById('cvSection');
    section.classList.toggle('show');
  }
  
  
document.addEventListener("DOMContentLoaded", () => {
  const forms = document.querySelectorAll(".json-form");

  forms.forEach(form => {
    form.addEventListener("submit", async e => {
      e.preventDefault();

      const endpoint = form.action || window.location.pathname;
      const feedbackId = form.getAttribute("data-feedback-id") || "form-feedback";
      let feedback = document.getElementById(feedbackId);

      if (!feedback) {
        feedback = document.createElement("div");
        feedback.id = feedbackId;
        feedback.style.marginTop = "1rem";
        form.appendChild(feedback);
      }

      feedback.textContent = "";
      feedback.style.color = "black";

      const formData = new FormData(form);
      const jsonData = {};

      formData.forEach((value, key) => {
        jsonData[key] = value;
      });

      const response = await fetch(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(jsonData),
      });

      const data = await response.json().catch(() => ({}));
      if (response.ok) {

        // our fun redirect
        if (data.redirect) {
          showRedirectCountdown(3, data.redirect); // 👈 this does the magic
          return;
        }

        feedback.style.color = "green";
        feedback.textContent = `✅ ${data.server || "Success!"}`;
        form.reset();
      } else {
        feedback.style.color = "red";
        feedback.textContent = `❌ ${data.error || "Something went wrong."}`;
      }
    });
  });
});


// our fun lil redirect
function showRedirectCountdown(seconds, targetUrl) {
  const messageEl = document.getElementById("redirect-message");
  messageEl.style.display = "block";

  let counter = seconds;

  const updateText = () => {
    messageEl.textContent = `✅ Login successful. You’ll be redirected in ${counter} second${counter === 1 ? "" : "s"}...`;
  };

  updateText();

  const interval = setInterval(() => {
    counter--;
    updateText();

    if (counter <= 0) {
      clearInterval(interval);
      window.location.href = targetUrl;
    }
  }, 1000);
}