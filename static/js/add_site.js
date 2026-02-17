const modal = document.getElementById('domainModal');
const closeBtn = document.querySelector('.close-button');
const form = document.getElementById('domainForm');
const title = document.getElementById('modalTitle');
const submitBtn = document.getElementById('submitButton');

// Open modal in Add mode
document.getElementById('openModalBtn').addEventListener('click', () => {
    modal.style.display = 'flex';
    openModal({
        mode: 'add',
    });
});

document.querySelectorAll('.edit-domain').forEach((btn) => {
    btn.addEventListener('click', (e) => {
        e.preventDefault();
        openModal({
            mode: 'edit',
            id: btn.dataset.id,
            domain: btn.dataset.domain,
            likes: btn.dataset.likes === 'true',
            comments: btn.dataset.comments === 'true',
        });
    });
});

function openModal({ mode, id, domain, likes, comments }) {
    modal.style.display = 'flex';
    document.body.style.overflow = 'hidden';

    if (mode === 'edit') {
        title.textContent = 'Edit Domain';
        submitBtn.textContent = 'Update Domain';
        form.action = `/edit_domain/${id}`;
        document.getElementById('domainId').value = id;
        document.getElementById('domainName').value = domain || '';
        document.getElementById('enableLike').checked = likes || false;
        document.getElementById('enableComments').checked = comments || false;
    } else {
        title.textContent = 'Add Domain';
        submitBtn.textContent = 'Add Domain';
        form.action = '/add_domain';
        form.reset();
    }
}

closeBtn.addEventListener('click', () => {
    modal.style.display = 'none';
    document.body.style.overflow = 'auto';
});

window.addEventListener('click', (e) => {
    if (e.target === modal) {
        modal.style.display = 'none';
        document.body.style.overflow = 'auto';
    }
});

document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('domainForm');
    const domainInput = document.getElementById('domainName');
    const errorMessage = document.getElementById('domainError');

    // Domain validation regex (no protocol, no paths)
    const domainRegex = /^(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(?::\d{1,5})?$/;

    // Auto-clean input: removes http(s)://, www., and trailing slashes
    function cleanDomainInput(value) {
        return value
            .trim()
            .toLowerCase()
            .replace(/^https?:\/\//, '') // remove http:// or https://
            .replace(/^www\./, '') // remove www.
            .replace(/\/.*$/, ''); // remove anything after first slash
    }

    function validateDomain(value) {
        const cleaned = cleanDomainInput(value);
        domainInput.value = cleaned; // auto-corrected value

        if (!cleaned) {
            errorMessage.textContent = 'Domain name is required.';
            return false;
        } else if (!domainRegex.test(cleaned)) {
            errorMessage.textContent =
                'Please enter a valid domain name (e.g., example.com).';
            return false;
        } else {
            errorMessage.textContent = '';
            return true;
        }
    }

    // Validate on input change
    domainInput.addEventListener('input', () => {
        validateDomain(domainInput.value);
    });

    // Validate before form submit
    form.addEventListener('submit', (e) => {
        if (!validateDomain(domainInput.value)) {
            e.preventDefault();
            domainInput.focus();
        }
    });
});
