const express = require('express');
const app = express();
const net = require('net');
const path = require('path');
const { DateTime } = require("luxon");

app.use(express.static('public'));

app.use((req, res, next) => {
  if (req.path === '/') {
    res.redirect('/create_link');
  } else {
    next();
  }
});

app.post('/create_link', (req, res) => {
  let inputValue = ``;

  req.on('data', (chunk) => {
    inputValue += chunk;
  });

  req.on('end', () => {
    const client = new net.Socket();
    client.connect(6379, '127.0.0.1', () => {
      const requestData = {
        method: 'POST',
        table: 'links',
        link: inputValue,
    };

    const requestString = JSON.stringify(requestData);
    client.write(requestString);
    });

    let responseData = '';

    client.on('data', (data) => {
      responseData += data.toString();
    });

    client.on('end', () => {
      const short_link = `http://31.28.27.213/${responseData}`
      res.end(short_link);
      client.end();
    });

    client.on('error', (error) => {
      console.error(`Ошибка при подключении к серверу: ${error}`);
      res.status(500).end('Internal Server Error');
    });
  });
});

app.get('/:value', (req, res) => {
  const value = req.params.value;
  const isExplicitRedirect = req.query.explicit === 'true';
  console.log(value)
  if (value === "favicon.ico") {
    res.status(204).end();
    return; 
  }
  if (value === "element not found") {
    res.status(204).end();
    return; 
  }

  if (value === "create_link"){
    const filePath = path.resolve("site", 'index.html');
    res.sendFile(filePath);
    return; 
  }

  if (value === "report"){
    const client = new net.Socket();
    client.connect(1333, '127.0.0.1', () => {
      const requestData = {
        method: 'GET',
        table: 'statistics',
        link: `1`,
        ip:   `2`,
        time_interval: `3`,
    };

      const requestString = JSON.stringify(requestData);
      client.write(requestString);
    });

    let responseData = ''; 

    client.on('data', (data) => {
      responseData += data.toString().trim(); 
    });
    
    client.on('end', () => {  
      client.end();
    });
    return;  
    }

  const client = new net.Socket();
  client.connect(6379, '127.0.0.1', () => {
    const requestData = {
      method: 'GET',
      table: 'links',
      link: value,
    };

    const requestString = JSON.stringify(requestData);
    client.write(requestString);
  });

  let responseData = ''; 

  client.on('data', (data) => {
    responseData += data.toString().trim(); 
  });

  client.on('end', () => {

    const client1 = new net.Socket();
    client1.connect(1333, '127.0.0.1', () => {

        const localIp = req.ip;
        const timeZone = "Asia/Novosibirsk";
        const currentDate = DateTime.now().setZone(timeZone);
        const formattedCurrentTime = currentDate.toFormat("HH:mm");

        const nextMinute = currentDate.plus({ minutes: 1 });
        const formattedNextMinuteTime = nextMinute.toFormat("HH:mm");

        const requestData = { 
            method: 'POST',
            table: 'links',
            link: value,
            ip: localIp,
            time_interval: `${formattedCurrentTime}-${formattedNextMinuteTime}`
        };

        const requestString = JSON.stringify(requestData);

        client1.write(requestString, () => {
            client1.end();
        });
    });

    client1.on('close', () => {
    });

    res.redirect(responseData);

    client.end();
});

  client.on('error', (error) => {
    console.error(`Ошибка при подключении к серверу: ${error}`);
    res.status(500).send('Internal Server Error');
  });
});

app.listen(8080, '127.0.0.1', () => {
  console.log('Сервер запущен на порту 80');
});

function sendDataToSecondServer(req, value, responseData) {
  const client1 = new net.Socket();
  client1.connect(1333, '127.0.0.1', () => {
    console.log('Connected to the second server');

    const localIp = req.ip;
    const timeZone = "Asia/Novosibirsk";
    const currentDate = DateTime.now().setZone(timeZone);
    const formattedCurrentTime = currentDate.toFormat("HH:mm");

    const nextMinute = currentDate.plus({ minutes: 1 });
    const formattedNextMinuteTime = nextMinute.toFormat("HH:mm");

    const requestData = {
      method: 'POST',
      table: 'links',
      link: value,
      ip: localIp,
      time_interval: `${formattedCurrentTime}-${formattedNextMinuteTime}`
    };

    const requestString = JSON.stringify(requestData);
    console.log('Sending request to the second server:', requestString);

    client1.write(requestString, () => {
      console.log('Request sent successfully to the second server');
      client1.end();
      console.log('Closed the connection to the second server');
    });
  });

  client1.on('close', () => {
    console.log('Connection to the second server closed');
  });
}