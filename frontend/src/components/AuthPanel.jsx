import { useState } from 'react';
import { api, extractToken } from '../services/api';

export default function AuthPanel({ token, onTokenChange, onNotify }) {
  const [mode, setMode] = useState('login');
  const [form, setForm] = useState({ name: '', email: '', password: '' });
  const [manualToken, setManualToken] = useState(token || '');
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e) {
    e.preventDefault();
    setLoading(true);
    try {
      const payload =
        mode === 'register'
          ? { name: form.name, email: form.email, password: form.password }
          : { email: form.email, password: form.password };

      const res = mode === 'register' ? await api.register(payload) : await api.login(payload);
      const foundToken = extractToken(res);

      if (foundToken) {
        onTokenChange(foundToken);
        setManualToken(foundToken);
      }
      onNotify(
        'success',
        mode === 'register'
          ? 'Registrasi berhasil!'
          : 'Login berhasil! Token JWT disimpan.'
      );
      setForm((prev) => ({ ...prev, password: '' }));
    } catch (err) {
      onNotify('error', err.message);
    } finally {
      setLoading(false);
    }
  }

  function handleSaveManualToken(e) {
    e.preventDefault();
    onTokenChange(manualToken.trim());
    onNotify('success', manualToken.trim() ? 'Token Bearer diperbarui.' : 'Token dikosongkan.');
  }

  function handleLogout() {
    onTokenChange('');
    setManualToken('');
    onNotify('info', 'Token dihapus dari browser.');
  }

  return (
    <section className="card">
      <div className="card-header">
        <div>
          <h2>1. Autentikasi (JWT / Bearer)</h2>
          <p className="muted">Login/Register atau tempel token JWT jika route diproteksi.</p>
        </div>
        <div className="tab-group">
          <button
            type="button"
            className={mode === 'login' ? 'btn btn-primary btn-sm' : 'btn btn-ghost btn-sm'}
            onClick={() => setMode('login')}
          >
            Login
          </button>
          <button
            type="button"
            className={mode === 'register' ? 'btn btn-primary btn-sm' : 'btn btn-ghost btn-sm'}
            onClick={() => setMode('register')}
          >
            Register
          </button>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="form-grid">
        {mode === 'register' && (
          <label>
            Nama Lengkap
            <input
              type="text"
              placeholder="Nama lengkap"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              required
            />
          </label>
        )}
        <label>
          Email
          <input
            type="email"
            placeholder="nama@email.com"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            required
          />
        </label>
        <label>
          Password
          <input
            type="password"
            placeholder="••••••••"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
            required
          />
        </label>
        <div className="form-actions">
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Memproses...' : mode === 'login' ? 'Masuk' : 'Daftar Akun'}
          </button>
          {token && (
            <button type="button" className="btn btn-danger" onClick={handleLogout}>
              Hapus Token
            </button>
          )}
        </div>
      </form>

      <form onSubmit={handleSaveManualToken} className="token-row">
        <input
          type="text"
          placeholder="Tempel token JWT manual..."
          value={manualToken}
          onChange={(e) => setManualToken(e.target.value)}
        />
        <button type="submit" className="btn btn-secondary btn-sm">
          Simpan Token
        </button>
      </form>

      <div className="status-strip">
        Status Token:{' '}
        <strong>{token ? `Aktif (${token.slice(0, 18)}...)` : 'Belum ada token'}</strong>
      </div>
    </section>
  );
}
