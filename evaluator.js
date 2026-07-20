// ~/cli/evaluator.js

/**
 * Evalúa un objeto de cambio proveniente de Gerrit.
 * @param {Object} cambio - Objeto JSON del cambio.
 * @returns {Object} - Objeto procesado con indicadores de criticidad.
 */
export const evaluarCambio = (cambio) => {
  return {
    id: cambio._number,
    resumen: cambio.subject,
    // Detecta palabras clave para priorización
    esCritico: /fix|security|critical/i.test(cambio.subject),
    url: `https://gerrit.wikimedia.org/r/c/${cambio._number}`
  };
};
