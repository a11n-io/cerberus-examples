import {clsx} from "clsx";
import {Button} from "react-bootstrap";

export default function Btn(props) {
    const { children, outline, accent, className, ...rest } = props;

    const classNames = clsx({
            btn: true,
            "btn-default": !outline,
            "btn-outline": outline,
            "btn-accent": accent
        },
        className
    );

    return (
        <Button className={classNames} {...rest}>
            {children}
        </Button>
    );
}