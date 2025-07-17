import { FC } from "react";

const About: FC = () => {
  return (
    <section id="about" className="py-16">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div>
            <h2 className="text-3xl md:text-4xl font-bold text-white mb-6">
              About Orvex
            </h2>
            <p className="text-lg text-gray-300 mb-6 leading-relaxed">
              We're passionate about providing the best Minecraft hosting
              experience around the world. Since our founding, we've been
              dedicated to delivering high-performance servers at unbeatable
              prices.
            </p>
            <p className="text-gray-400 mb-6 leading-relaxed">
              Our mission is simple: make quality Minecraft hosting accessible
              to everyone. Whether you're running a small server with friends or
              managing a large public community, we have the perfect solution
              for you.
            </p>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <FeatureCard
              icon={<i className="fas fa-users text-white" />}
              title="Community Focused"
              description="We understand the Minecraft community and build our services around your needs."
            />
            <FeatureCard
              icon={<i className="fas fa-award text-white" />}
              title="Quality First"
              description="We never compromise on quality, ensuring your server runs perfectly every time."
            />
            <FeatureCard
              icon={<i className="fas fa-clock text-white" />}
              title="Instant Support"
              description="Our expert team is available 24/7 to help you with any questions or issues."
            />
            <FeatureCard
              icon={<i className="fas fa-map-pin text-white" />}
              title="EU/US Presence"
              description="Strategic locations in the EU and US ensure low latency for all your players."
            />
          </div>
        </div>
      </div>
    </section>
  );
};

type FeatureCardProps = {
  icon: React.ReactNode;
  title: string;
  description: string;
};

const FeatureCard: FC<FeatureCardProps> = ({ icon, title, description }) => (
  <div className="bg-black/30 p-6 rounded-xl border border-white/10 hover:border-white/20 transition-all duration-300">
    <div className="bg-white/10 w-12 h-12 rounded-lg flex items-center justify-center mb-4">
      {icon}
    </div>
    <h3 className="text-lg font-semibold text-white mb-2">{title}</h3>
    <p className="text-gray-400 text-sm">{description}</p>
  </div>
);

export default About;
