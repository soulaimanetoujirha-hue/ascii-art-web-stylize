// Small touches on top of a page that already works without JavaScript:
// live feedback while typing, copy and download, and fitting the art to its box.

const form = document.getElementById("form");
const input = document.getElementById("text");
const statusLine = document.getElementById("status");
const statusText = document.getElementById("status-text");
const art = document.getElementById("art");

function say(message, state = "") {
    statusText.textContent = message; 
    statusLine.dataset.state = state; 
}

// Same rule as the server: printable ASCII (space to ~) plus line breaks.
function checkInput() {
    const bad = [...input.value].find((c) => c !== "\n" && c !== "\r" && (c < " " || c > "~"));

    if (bad) {
        input.setCustomValidity(`"${bad}" can't be drawn. Use letters, numbers, spaces and common symbols.`);
        input.classList.add("is-invalid");
        say(`"${bad}" can't be drawn, printable ASCII only`, "error");
        return;
    }

    input.setCustomValidity("");
    input.classList.remove("is-invalid");
    const count = input.value.replace(/\r?\n/g, "").length;
    const font = form.elements.banner.value;
    say(count ? `${count} character${count === 1 ? "" : "s"} · ${font}` : "Ready when you are");
}

// Shrink the font until the art fits its box (down to 7px, then it scrolls).
function fitArt() {
    const box = art.parentElement;
    let size = 16;
    art.style.fontSize = size + "px";
    while (size > 7 && (art.offsetWidth > box.clientWidth || art.offsetHeight > box.clientHeight)) {
        size -= 0.5;
        art.style.fontSize = size + "px";
    }
}

if (form) {
    input.addEventListener("input", checkInput);
    form.addEventListener("change", checkInput);
    form.addEventListener("submit", () => say("Drawing your art…", "busy"));

    // Ctrl + Enter (or Cmd + Enter) generates straight from the text box.
    input.addEventListener("keydown", (event) => {
        if (event.key === "Enter" && (event.ctrlKey || event.metaKey)) {
            form.requestSubmit();
        }
    });
}

if (art) {
    fitArt();
    window.addEventListener("resize", fitArt);

    // On phones the result sits below the form, so bring it into view.
    if (form && window.matchMedia("(max-width: 860px)").matches) {
        art.scrollIntoView({ block: "center" });
    }
}

document.getElementById("copy")?.addEventListener("click", async () => {
    try {
        await navigator.clipboard.writeText(art.textContent);
        say("Copied to your clipboard");
    } catch {
        // Some browsers block clipboard access: select the art so Ctrl + C still works.
        getSelection().selectAllChildren(art);
        say("Art selected, press Ctrl + C to copy");
    }
});

document.getElementById("download")?.addEventListener("click", () => {
    const link = document.createElement("a");
    link.href = URL.createObjectURL(new Blob([art.textContent], { type: "text/plain" }));
    link.download = "ascii-art.txt";
    link.click();
    setTimeout(() => URL.revokeObjectURL(link.href), 1000);
    say("Saved as ascii-art.txt");
});
