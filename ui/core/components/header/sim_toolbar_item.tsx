import clsx from 'clsx';

const Element = ({ children, buttonClassName, linkRef, ...props }: SimToolbarItemProps) =>
	'href' in props && props.href ? (
		<a ref={linkRef as JSX.HTMLElementProps<'a'>['ref']} href={props.href} target="_blank" className={clsx(buttonClassName)}>
			{children}
		</a>
	) : (
		<button type="button" ref={linkRef as JSX.HTMLElementProps<'button'>['ref']} className={clsx(buttonClassName)}>
			{children}
		</button>
	);

export type SimToolbarItemProps = {
	icon?: string;
	iconRef?: JSX.HTMLElementProps<'i'>['ref'];
	href?: string;
	linkRef?: JSX.HTMLElementProps<'button'>['ref'] | JSX.HTMLElementProps<'a'>['ref'];
	buttonClassName?: string;
} & Pick<JSX.HTMLElementProps<'div'>, 'ref' | 'children' | 'className'>;

export const SimToolbarItem = ({ ref, linkRef, iconRef, icon, className, children, ...props }: SimToolbarItemProps) => (
	<div ref={ref} className="sim-toolbar-item">
		<Element linkRef={linkRef} className={clsx(className)} {...props}>
			{icon && <i ref={iconRef} className={icon} />}
			{children}
		</Element>
	</div>
);
