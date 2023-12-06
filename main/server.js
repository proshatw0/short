const express = require('express');
const app = express();
const net = require('net');
const path = require('path');
const { DateTime } = require("luxon");
const bodyParser = require('body-parser');

app.use(express.static('public'));
app.use(bodyParser.json());

app.use((req, res, next) => {
  if (req.path === '/') {
    res.redirect('/create_link');
  } else {
    next();
  }
});

app.post('/create_report', (req, res) => {
  const { variable1, variable2, variable3 } = req.body;

  const client = new net.Socket();
  client.connect(1333, '127.0.0.1', () => {
    const requestData = {
      method: 'GET',
      table: 'statistics',
      link: variable2,
      ip: variable3,
      time_interval: variable1,
    };
    const requestString = JSON.stringify(requestData);
    client.write(requestString);
  });
  
  let responseData = '';
  
  client.on('data', (data) => {
    responseData += data.toString();
  });
  
  client.on('end', () => {
    try {
      const responseObject = JSON.parse(responseData);
  
      res.json(responseObject);
    } catch (error) {
      console.error(`Ошибка при разборе JSON: ${error}`);
      res.status(500).end('Internal Server Error');
    } finally {
      client.end();
    }
  });
  client.on('error', (error) => {
    console.error(`Ошибка при подключении к серверу: ${error}`);
    res.status(500).end('Internal Server Error');
  });
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
      const short_link = `http://31.28.27.213:3000/${responseData}`
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

  if (value === "create_report"){
    const filePath = path.resolve("site", 'report.html');
    res.sendFile(filePath);
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

app.listen(3000, '31.28.27.213', () => {
  console.log('Сервер запущен на порту 3000');
});