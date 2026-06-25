import https from 'https';

// PARÁMETROS DE ARQUITECTURA (Espejo Público de Wikimedia)
const GERRIT_HOST = 'gerrit.wikimedia.org'; 
const PROYECTO = 'mediawiki/core'; // El proyecto principal de Wikipedia
const RAMA = 'master';             // Rama de integración principal

// Consulta REST nativa: Cambios abiertos en mediawiki/core para la rama master
const PATH_CONSULTA = `/r/changes/?q=status:open+project:${PROYECTO}+branch:${RAMA}&n=5`;

const options = {
  hostname: GERRIT_HOST,
  path: PATH_CONSULTA,
  method: 'GET',
  headers: {
    'User-Agent': 'Diamond-DevOps-Automator/1.0',
    'Accept': 'application/json'
  }
};

console.log(`🚀 Diamond-Automator | Consultando Gerrit Público: ${GERRIT_HOST}`);
console.log(`🔍 Filtrando Proyecto: ${PROYECTO} | Rama: ${RAMA}...\n`);

const req = https.request(options, (res) => {
  let rawData = '';

  res.on('data', (chunk) => { rawData += chunk; });

  res.on('end', () => {
    try {
      // PILAR CRÍTICO: Remoción obligatoria del prefijo anti-XSS de Gerrit )]}'\n
      const cleanData = rawData.replace(/^\)\]\}'\n/, '');
      const cambios = JSON.parse(cleanData);

      console.log(`✅ Conexión exitosa. Cambios pendientes detectados: ${cambios.length}\n`);
      
      cambios.forEach((cambio) => {
        console.log(`[-] [ID: ${cambio._number}] ${cambio.subject}`);
        console.log(`    Estado: ${cambio.status} | Rama: ${cambio.branch}`);
        console.log(`    URL de Revisión: https://${GERRIT_HOST}/r/c/${cambio._number}\n`);
      });

    } catch (error) {
      console.error('❌ Error al procesar el buffer JSON:', error.message);
    }
  });
});

req.on('error', (e) => {
  console.error(`❌ Fallo en la comunicación HTTPS directa: ${e.message}`);
});

req.end();
