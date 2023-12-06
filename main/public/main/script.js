function copyLink() {
  const button = document.getElementById("Link");
  const buttonText = button.textContent;
  const value = buttonText

  const tempInput = document.createElement("input");
  tempInput.value = value;

  document.body.appendChild(tempInput);

  tempInput.select();

  try {
    document.execCommand("copy");
  } catch (err) {
    console.error("Не удалось скопировать текст: ", err);
  }
  document.body.removeChild(tempInput);
}

function create_link() {
    const linkInput = document.getElementById('linkInput');
    const inputValue = linkInput.value;
  
    fetch('/create_link', {
      method: 'POST',
      body: inputValue, 
      headers: {
        'Content-Type': 'text/plain' 
      }
    })
    .then(response => response.text()) 
    .then(data => {
      container = document.querySelector(".Input")
      stroke = document.querySelector(".stroke")
      if (data.includes("31.28.27.213:3000")) {
          const errorElement = document.getElementById('Error');
        if (errorElement) {
          errorElement.remove();
        }
        const testTextLinksElement = document.getElementById('textLinks');
        if (testTextLinksElement){
          testTextLinksElement.remove();
        }
        var textLinks = document.createElement('h2');
        textLinks.id = `textLinks`
        stroke.appendChild(textLinks)

        const testLinksElement = document.getElementById('Link');
        if (testLinksElement){
          testLinksElement.remove();
        }
        var linkElement = document.createElement('button');
          linkElement.className = 'Link';
          linkElement.id = 'Link';
          linkElement.textContent = data;
          linkElement.addEventListener('click', copyLink);
          container.appendChild(linkElement)
      } else {
        const textLinks = document.getElementById('textLinks');
        if (textLinks) {
          textLinks.remove();
        }

        const linkElement = document.getElementById('Link');
        if (linkElement) {
          linkElement.remove();
        }

        const testErrorElement = document.getElementById('Error');
        if (testErrorElement){
          testErrorElement.remove();
        }
        var errorElement = document.createElement('p');
        errorElement.className = 'Error';
        errorElement.id = 'Error';
        errorElement.textContent = `Что-то пошло не так: ${data}`;
        container.appendChild(errorElement)
      }
    })
    .catch(error => {
      console.error('Ошибка:', error);
    });
}
