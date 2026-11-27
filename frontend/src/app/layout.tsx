import type { Metadata } from "next";
import { IBM_Plex_Sans } from "next/font/google";
import AppShell from "@/components/layout/AppShell";
import "./globals.css";

const ibmPlexSans = IBM_Plex_Sans({
  subsets: ["latin"],
  weight: ["400", "500", "600"],
  variable: "--font-ibm-plex",
});

export const metadata: Metadata = {
  title: "RDO · Operações de Campo",
  description: "Plataforma de gestão de obras, serviços e operações de campo.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="pt-BR" className={ibmPlexSans.variable}>
      <body>
        <AppShell>{children}</AppShell>
      </body>
    </html>
  );
}