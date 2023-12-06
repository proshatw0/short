function validateInput() {
    var linkValue = document.getElementById('link').value;
    var timeValue = document.getElementById('ip').value;
    var ipValue = document.getElementById('time').value;

    var uniqueValues = [];

    if (linkValue.trim() !== '' && parseFloat(linkValue) !== 0) {
        if (!uniqueValues.includes(linkValue)) {
            uniqueValues.push(linkValue);
        } else {
            return false
        }
    }

    if (timeValue.trim() !== '' && parseFloat(timeValue) !== 0) {
        if (!uniqueValues.includes(timeValue)) {
            uniqueValues.push(timeValue);
        } else {
            return false
        }
    } 

    if (ipValue.trim() !== '' && parseFloat(ipValue) !== 0) {
        if (!uniqueValues.includes(ipValue)) {
            uniqueValues.push(ipValue);
        } else {
            return false
        }
    }

    if ((uniqueValues.includes('2') || uniqueValues.includes('3')) && !uniqueValues.includes('1')){
        return false
    }

    if ((uniqueValues.includes('1') && uniqueValues.includes('3')) && !uniqueValues.includes('2')){
        return false
    }
    if ((uniqueValues.includes('2') && uniqueValues.includes('3')) && !uniqueValues.includes('1')){
        return false
    }
    if (uniqueValues.length == 0){
        return false
    }

    return true
}

function create_report() {
    var linkValid = validateInput();

    if (!linkValid) {
        alert("Значения в полях ввода должны быть уникальными и иметь логическую последовательность.");
        return;
    }

    var linkValue = document.getElementById('link').value;
    var timeValue = document.getElementById('ip').value;
    var ipValue = document.getElementById('time').value;

    const data = {
        variable1: ipValue,
        variable2: linkValue,
        variable3: timeValue
    };
    
    const jsonString = JSON.stringify(data);
  
    fetch('/create_report', {
        method: 'POST',
        body: jsonString,
        headers: {
        'Content-Type': 'application/json'
    }
    })
    .then(response => response.json())
    .then(data => {
        function sortByPid(data, pid = 0, depth = 0) {
            const result = [];
            
            for (const item of data) {
              if (item.pid === pid) {
                item.depth = depth;
                result.push(item);
                result.push(...sortByPid(data, item.id, depth + 1));
              }
            }
          
            return result;
          }

        const sortedData = sortByPid(data.Table);

        const tableContainer = document.getElementById('TableBody');
        tableContainer.innerHTML = '';

        sortedData.forEach(item => {
            const row = document.createElement('tr');
        
            const cell1 = document.createElement('td');
        
            if (item.pid !== 0) {
                if (data.Table.find(element => element.id === item.pid && element.pid !== 0)) {
                    cell1.classList.add('double-padding');
                } else {
                    cell1.classList.add('with-padding');
                }
            }
        
            cell1.textContent = item.link || item.time_interval || item.ip;
        
            const cell2 = document.createElement('td');
            cell2.textContent = item.cout;
        
            row.appendChild(cell1);
            row.appendChild(cell2);
        
            tableContainer.appendChild(row);
        });
    })
    .catch(error => {
      console.error('Ошибка:', error);
    });
}