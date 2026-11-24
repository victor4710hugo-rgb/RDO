export default function Topbar({
  title,
  subtitle,
}: {
  title: string;
  subtitle?: string;
}) {
  return (
    <header className="topbar">
      <div>
        <h1>{title}</h1>
        {subtitle && <small>{subtitle}</small>}
      </div>
      <div className="avatar" aria-label="Perfil">
        L
      </div>
    </header>
  );
}