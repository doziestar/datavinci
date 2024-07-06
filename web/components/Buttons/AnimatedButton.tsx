import { motion } from "framer-motion";
import { Button } from "../ui/button";
import { useAnimations } from "@/hooks/animation/useAnimation";

const AnimatedButton: React.FC<React.ComponentProps<typeof Button>> = ({
  children,
  ...props
}) => {
  const { splash, jelly, smoothAppear, wobble } = useAnimations();

  return (
    <motion.div whileTap="tap" variants={jelly}>
      <Button {...props}>{children}</Button>
    </motion.div>
  );
};

export { AnimatedButton };
