"use client";

const About: React.FC = () => {
  const stats = [
    {
      number: "99.9%",
      label: "Uptime Guarantee",
      icon: "fas fa-shield-alt",
    },
    {
      number: "24/7",
      label: "Expert Support",
      icon: "fas fa-headset",
    },
    {
      number: "5min",
      label: "Server Deployment",
      icon: "fas fa-rocket",
    },
    {
      number: "2+",
      label: "Global Locations",
      icon: "fas fa-globe",
    },
  ];

  const values = [
    {
      title: "Performance First",
      description: "We use cutting-edge hardware and optimize every aspect of our infrastructure to deliver the best gaming experience possible.",
      icon: "fas fa-tachometer-alt",
    },
    {
      title: "Reliability",
      description: "Our enterprise-grade infrastructure ensures your server stays online when your community needs it most.",
      icon: "fas fa-server",
    },
    {
      title: "Simplicity",
      description: "From deployment to management, we've made everything as simple as possible so you can focus on what matters - gaming.",
      icon: "fas fa-magic",
    },
  ];

  return (
    <section id="about" className="py-20 px-8 bg-neutral-950 text-white relative">
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl font-semibold mb-4 text-transparent bg-clip-text bg-gradient-to-r from-neutral-400 via-white to-neutral-300">
            About Orvex
          </h2>
          <p className="text-lg font-normal text-neutral-400 max-w-3xl mx-auto">
            Founded by gamers, for gamers. We understand what it takes to run the perfect Minecraft server 
            because we've been there ourselves. Our mission is simple: provide the most reliable, fast, 
            and affordable Minecraft hosting service in the industry.
          </p>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-6 mb-20">
          {stats.map((stat, index) => (
            <div
              key={index}
              className="p-6 bg-neutral-900/50 backdrop-blur-sm border border-white/10 text-center rounded-3xl"
            >
              <div className="flex flex-col items-center">
                <div className="w-12 h-12 bg-white/10 rounded-lg flex items-center justify-center mb-4">
                  <i className={`${stat.icon} text-violet-400 text-xl`}></i>
                </div>
                <div className="text-2xl font-semibold text-white mb-1">{stat.number}</div>
                <div className="text-neutral-400 text-sm font-normal">{stat.label}</div>
              </div>
            </div>
          ))}
        </div>

        <div className="mb-16">
          <h3 className="text-3xl font-semibold mb-8 text-center text-transparent bg-clip-text bg-gradient-to-r from-neutral-400 via-white to-neutral-300">
            Our Values
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {values.map((value, index) => (
              <div
                key={index}
                className="p-6 bg-neutral-900/50 backdrop-blur-sm border border-white/10 rounded-3xl"
              >
                <div className="flex flex-col h-full">
                  <div className="w-10 h-10 bg-white/10 rounded-lg flex items-center justify-center mb-4">
                    <i className={`${value.icon} text-violet-400 text-lg`}></i>
                  </div>
                  <h4 className="text-lg font-normal mb-2 text-white">
                    {value.title}
                  </h4>
                  <p className="text-neutral-400 text-sm font-normal leading-relaxed">
                    {value.description}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="text-center">
          <div className="bg-gradient-to-r from-violet-600/20 to-purple-600/20 rounded-2xl p-8 border border-white/10">
            <h3 className="text-2xl font-semibold mb-4 text-white">
              Why Choose Orvex?
            </h3>
            <p className="text-neutral-300 text-lg max-w-4xl mx-auto leading-relaxed">
              We're not just another hosting company. We're a team of passionate gamers and developers 
              who have spent years optimizing server performance, building intuitive control panels, 
              and providing the kind of support we'd want for our own servers. When you choose Orvex, 
              you're choosing a partner who truly understands gaming.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};

export default About;