const services = [
  {
    icon: <i className="fas fa-server text-white" />,
    title: "High-Performance Hardware",
    description:
      "AMD Ryzen CPUs, DDR4 RAM, and NVMe storage for smooth performance.",
  },
  {
    icon: <i className="fas fa-shield-alt text-white" />,
    title: "DDoS Protection",
    description: "Enterprise-grade protection keeps your server online 24/7.",
  },
  {
    icon: <i className="fas fa-bolt text-white" />,
    title: "Instant Setup",
    description:
      "Your server is ready in under 5 minutes with automated deployment.",
  },
  {
    icon: <i className="fas fa-cog text-white" />,
    title: "Full Control Panel",
    description:
      "Manage your server with our intuitive control panel interface.",
  },
  {
    icon: <i className="fas fa-headphones text-white" />,
    title: "24/7 Expert Support",
    description:
      "Our Minecraft experts are available around the clock to help.",
  },
  {
    icon: <i className="fas fa-globe text-white" />,
    title: "EU/US Presence",
    description: "Low latency with strategically located EU/US data centers.",
  },
];

const Services: React.FC = () => {
  return (
    <section id="services" className="py-16">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
            Why Choose Orvex?
          </h2>
          <p className="text-lg text-gray-300 max-w-2xl mx-auto">
            Everything you need for a successful Minecraft server with unmatched
            performance.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {services.map(({ icon, title, description }, idx: number) => (
            <div
              key={idx}
              className="bg-black/30 p-6 rounded-xl border border-white/10 hover:border-white/20 transition-all duration-300 hover:transform hover:scale-105"
            >
              <div className="bg-white/10 w-12 h-12 rounded-lg flex items-center justify-center mb-4">
                {icon}
              </div>
              <h3 className="text-lg font-semibold text-white mb-3">{title}</h3>
              <p className="text-gray-400 text-sm">{description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default Services;
