import { evaluarCambio } from './evaluator.js';

const URL = 'https://gerrit.wikimedia.org/r/changes/?q=status:open+project:mediawiki/core+branch:master&n=5';

async function consultarCambios() {
  try {
    const res = await fetch(URL, {
      headers: {
        'User-Agent': 'Diamond-DevOps-Automator/1.0',
        'Accept': 'application/json'
      }
    });

    if (!res.ok) {
      throw new Error(`HTTP Error status: ${res.status}`);
    }

    const rawText = await res.text();
    // Limpieza robusta del prefijo XSS/Security de Gerrit
    const cleanData = rawText.replace(/^\)\]\}'/, '').trim();

    const cambios = JSON.parse(cleanData);
    const resultados = cambios.map(evaluarCambio);

    console.table(resultados);
  } catch (error) {
    console.error('❌ Error en el proceso:', error.message);
  }
}

consultarCambios();
