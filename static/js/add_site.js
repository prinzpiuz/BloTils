document.addEventListener("DOMContentLoaded", function () {
    const modal = document.getElementById("add_domain");
    const openModalBtn = document.getElementById("openModalBtn");
    const closeModalBtn = document.querySelector(".close-button");
    const domainForm = document.getElementById("domainForm");
    const domainNameInput = document.getElementById("domainName");

    openModalBtn.addEventListener("click", function () {
         modal.style.display = 'flex';
    });

    closeModalBtn.addEventListener("click", function () {
        modal.style.display = "none";
    });

    window.addEventListener("click", function (event) {
        if (event.target === modal) {
            modal.style.display = "none";
        }
    });

    domainForm.addEventListener("submit", function (event) {
        const domainValue = domainNameInput.value.trim();
        if (domainValue === "") {
            alert("Please enter a domain name.");
            event.preventDefault(); // Prevent form submission
            return;
        }
        const cleanedDomain = domainValue.replace(/^https?:\/\//i, "");
        domainNameInput.value = cleanedDomain;

    });
});
