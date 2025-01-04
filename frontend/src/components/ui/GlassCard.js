const GlassCard = ({ children, className = "" }) => (
    <div
        className={`bg-white bg-opacity-5 backdrop-filter backdrop-blur-lg rounded-lg p-4 ${className}`}
    >
        {children}
    </div>
);
export default GlassCard;
