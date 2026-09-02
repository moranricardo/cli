const PATRON_CRITICO = /(fix|security|critical|urgent|hotfix|vuln|CVE-\d+)/i;

export const evaluarCambio = (cambio, { baseUrl = 'https://gerrit.wikimedia.org/r' } = {}) => {
  if (!cambio || typeof cambio._number !== 'number') {
    throw new TypeError('cambio inválido: requiere _number');
  }
  const subject = String(cambio.subject ?? '');
  return {
    id: cambio._number,
    resumen: subject.trim(),
    esCritico: PATRON_CRITICO.test(subject),
    prioridad: PATRON_CRITICO.test(subject) ? 'alta' : 'normal',
    url: `${baseUrl}/c/${cambio._number}`,
    owner: cambio.owner?.name ?? 'desconocido',
  };
};
