import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Rewrite /api/* para o backend Go — evita problemas de CORS em produção
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080"}/api/:path*`,
      },
      {
        source: "/webhooks/:path*",
        destination: `${process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080"}/webhooks/:path*`,
      },
    ];
  },

  // Domínios permitidos para next/image
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "*.supabase.co",
        pathname: "/storage/v1/object/public/**",
      },
    ],
  },

  // Headers de segurança
  async headers() {
    return [
      {
        source: "/(.*)",
        headers: [
          { key: "X-Content-Type-Options", value: "nosniff" },
          { key: "X-Frame-Options", value: "DENY" },
          { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
        ],
      },
    ];
  },
};

export default nextConfig;
