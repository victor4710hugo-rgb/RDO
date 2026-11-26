import sqlite3, sys

con = sqlite3.connect(f"file:{sys.argv[1]}?mode=ro", uri=True)
q = lambda s: con.execute(s).fetchall()

tabelas = [r[0] for r in q(
    "select name from sqlite_master where type='table' "
    "and name not like 'sqlite_%' order by name")]
print("TABELAS:", len(tabelas), tabelas)

if "schema_migrations" in tabelas:
    print("MIGRATIONS:", q("select * from schema_migrations"))

alvo = ["users", "workers", "organizations", "teams", "team_members",
        "roles", "permissions", "role_permissions", "user_roles",
        "project_memberships", "devices", "audit_logs"]
for t in alvo:
    if t not in tabelas:
        print(f"\n== {t}: NAO EXISTE")
        continue
    n = q(f"select count(*) from {t}")[0][0]
    print(f"\n== {t} ({n} linhas)")
    for c in q(f"pragma table_info({t})"):
        print("  ", c)
    for f in q(f"pragma foreign_key_list({t})"):
        print("  FK", f)
    for i in q(f"pragma index_list({t})"):
        print("  IDX", i)

print("\nforeign_key_check:", q("pragma foreign_key_check"))