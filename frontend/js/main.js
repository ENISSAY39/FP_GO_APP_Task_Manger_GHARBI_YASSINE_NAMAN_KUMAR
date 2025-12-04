// not used right now but useful if you want central init later
import { validate } from './lib/api.js';
export async function initApp() {
  try {
    const v = await validate();
    return v;
  } catch (e) {
    return null;
  }
}
