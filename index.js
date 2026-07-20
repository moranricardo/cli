import https from 'https';
import { evaluarCambio } from './evaluator.js';

const options = {
  hostname: 'gerrit.wikimedia.org',
  path: '/r/changes/?q=status:open+project:mediawiki/core+branch:master&n=5',
  method: 'GET',
  headers: { 
    'User-Agent': 'Diamond-DevOps-Automator/1.0', 
    'Accept': 'application/json' 
  }
};

const req = https.request(options, (res) => {
  let rawData = '';
  res.on('data', (chunk) => { rawData += chunk; });
  res.on('end', () => {
    // Limpieza agresiva del prefijo de seguridad de Gerrit
    const cleanData = rawData.replace(/^\)\]\}'\n/, '');
    
    try {
      const cambios = JSON.parse(cleanData);
      const resultados = cambios.map(evaluarCambio);
      console.table(resultados);
    } catch (e) {
      console.error('❌ Error al procesar JSON. Respuesta cruda:');
      console.error(cleanData.substring(0, 100));
    }
  });
});

req.on('error', (e) => console.error(`❌ Fallo de red: ${e.message}`));
req.end();

