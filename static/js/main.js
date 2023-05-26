var regexPattern;
var textArea;
var errorContainer;


document.addEventListener("DOMContentLoaded", function() {
    highlightGroups();
    regexPattern = document.getElementById("regexPattern");
    textArea = document.getElementById("textValue");
    errorContainer = document.getElementById("errorContainer");

    textArea.addEventListener("input", handleRegexPatternInput);
    regexPattern.addEventListener("input", handleRegexPatternInput);


    textArea.addEventListener("click", function() {
        textArea.focus();
    });
});

function highlightGroups() {
    const regexInput = document.getElementById('test');
    const pattern = regexInput.innerText;

    const openingBracket = '(';
    const closingBracket = ')';

    const stack = [];
    let groupCounter = 0;

    let output = '';

    for (let i = 0; i < pattern.length; i++) {
        const char = pattern[i];

        if (char === openingBracket) {
            groupCounter++;
            stack.push(groupCounter);
            output += `<span class="group-${groupCounter} bold">${char}</span>`;
        } else if (char === closingBracket) {
            const groupNumber = stack.pop();
            output += `<span class="group-${groupNumber} bold">${char}</span>`;
        } else {
            output += char;
        }
    }

    regexInput.innerHTML = output;
}
