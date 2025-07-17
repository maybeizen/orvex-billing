import { FC } from "react";

const Contact: FC = () => {
  return (
    <section id="contact" className="py-16">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
            Get in Touch
          </h2>
          <p className="text-lg text-gray-300 max-w-2xl mx-auto">
            Have questions about our services? Need help with your server? We're
            here to help!
          </p>
        </div>

        <div className="max-w-3xl mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <ContactCard
              icon={<i className="fas fa-envelope text-2xl text-white" />}
              title="Email Support"
              subtitle="Get help via email"
              linkText="support@orvex.cc"
              href="mailto:support@orvex.cc"
            />

            <ContactCard
              icon={<i className="fab fa-discord text-2xl text-white" />}
              title="Discord Community"
              subtitle="Join our community"
              linkText="discord.gg/orvex"
              href="https://discord.gg/orvex"
            />

            <ContactCard
              icon={<i className="fas fa-map-pin text-2xl text-white" />}
              title="Location"
              subtitle="Based in"
              extraText="New York, USA"
            />
          </div>
        </div>
      </div>
    </section>
  );
};

type ContactCardProps = {
  icon: React.ReactNode;
  title: string;
  subtitle: string;
  linkText?: string;
  href?: string;
  extraText?: string;
};

const ContactCard: FC<ContactCardProps> = ({
  icon,
  title,
  subtitle,
  linkText,
  href,
  extraText,
}) => (
  <div className="bg-zinc-950 p-6 rounded-xl border border-white/10 hover:border-white/20 transition-all duration-300 text-center">
    <div className="bg-white/10 w-14 h-14 rounded-xl flex items-center justify-center mx-auto mb-4">
      {icon}
    </div>
    <h4 className="text-white font-semibold mb-2">{title}</h4>
    <p className="text-gray-400 mb-3 text-sm">{subtitle}</p>
    {href && linkText ? (
      <a
        href={href}
        className="text-white hover:text-gray-300 font-medium text-sm"
      >
        {linkText}
      </a>
    ) : extraText ? (
      <p className="text-white font-medium text-sm">{extraText}</p>
    ) : null}
  </div>
);

export default Contact;
