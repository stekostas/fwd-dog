function copyLink() {
    const copyText = document.getElementById("short-link");

    copyText.select();
    copyText.setSelectionRange(0, 99999); /*For mobile devices*/

    document.execCommand("copy");

    const copyBtn = document.getElementById("copy-btn");
    const original = copyBtn.innerText;

    copyBtn.innerText = "COPIED";
    copyBtn.setAttribute("disabled", "disabled")

    setTimeout(function () {
        copyBtn.innerText = original;
        copyBtn.removeAttribute("disabled")
    }, 500);
}
