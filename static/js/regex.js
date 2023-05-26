
function handleRegexPatternInput() {
    const patternValue = regexPattern.value;
    // const textAreaValue = textArea.textContent;
    const textAreaValue = textArea.innerText;


    if (regexPattern !== "") {
        performRegexRequest(patternValue, textAreaValue);
    }
}

function performRegexRequest(regexPattern, input) {
    const data = new URLSearchParams();
    data.append("regexPattern", regexPattern);
    data.append("textValue", input);

    return fetch("/regex", {
        method: "POST",
        body: data
    })
        .then(function (response) {
            if (!response.ok) {
                throw new Error(`Request failed with status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            if (data.Error) {
                showErrorText(data.Error);
                return;
            }

            clearErrorText();

            updateInnerHTML(data.HTML);
        })
        .catch(error => {
            throw new Error(`Request failed with status: ${error}`);
        });

}

function showErrorText(error) {
    errorContainer.textContent = error;
}

function clearErrorText() {
    errorContainer.textContent = "";
}

function updateTextArea(html) {
    textArea.innerHTML = html;
}

function isTextAreaFocused() {
    return document.activeElement === textArea;
}

function updateInnerHTML(html) {
    let startPosition;
    const updateCursor = isTextAreaFocused();

    if (updateCursor) startPosition = getPosition();

    const currentTextContent = textArea.textContent;
    textArea.innerHTML = html;

    if (!updateCursor || !startPosition)
        return;

    const diff = extractTextFromHTML(html).length - currentTextContent.length;
    const newCursorPosition = startPosition + diff;

    restoreCursor(newCursorPosition);
}

function restoreCursor(position) {
    const range = document.createRange();
    const sel = window.getSelection();
    range.setStart(textArea, position);
    range.collapse(true);
    sel.removeAllRanges();
    sel.addRange(range);
}

function extractTextFromHTML(html) {
    const tempElement = document.createElement('div');
    tempElement.innerHTML = html;
    return tempElement.textContent;
}

function getPosition() {
    let position = null;
    if (!window.getSelection) {
        return null;
    }

    const sel = window.getSelection();
    if (sel.rangeCount === 0) {
        return null
    }

    const range = sel.getRangeAt(0);
    const clonedRange = range.cloneRange();
    clonedRange.selectNodeContents(textArea);
    clonedRange.setEnd(range.endContainer, range.endOffset);
    position = clonedRange.toString().length;
    return position;
}
