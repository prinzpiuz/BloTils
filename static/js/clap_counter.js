
function copyCode(btn) {
    const codeBlock = btn.closest('.code-block');
    const code = codeBlock.querySelector('code').textContent;

    navigator.clipboard
        .writeText(code)
        .then(() => {
            btn.classList.add('copied');
            btn.querySelector('span').textContent = 'Copied!';

            setTimeout(() => {
                btn.classList.remove('copied');
                btn.querySelector('span').textContent = 'Copy';
            }, 2000);
        })
        .catch((err) => {
            console.error('Failed to copy:', err);
        });
}
