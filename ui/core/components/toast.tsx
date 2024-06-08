import { Button, Icon, IconProps } from '@wowsims/ui';
import { Toast as BootstrapToast } from 'bootstrap';
import clsx from 'clsx';

type ToastOptions = {
	title?: string;
	variant: 'info' | 'success' | 'error' | 'warning';
	body: string | Element;
	autoShow?: boolean;
	canClose?: boolean;
	container?: Element;
	additionalClasses?: string[];
} & Partial<BootstrapToast.Options>;

class Toast {
	readonly element: HTMLElement;
	private container: Element;
	private title: ToastOptions['title'];
	private body: ToastOptions['body'];
	private variant: ToastOptions['variant'];
	private canClose: ToastOptions['canClose'];
	private additionalClasses: ToastOptions['additionalClasses'];

	public instance;
	constructor(options: ToastOptions) {
		const { title, variant, autoShow = true, canClose = true, body, additionalClasses, container, ...bootstrapOptions } = options || {};
		this.container = container || document.getElementById('toastContainer')!;
		this.additionalClasses = additionalClasses;
		this.title = title || 'WowSims';
		this.variant = variant || 'info';
		this.body = body;
		this.canClose = canClose;

		this.element = this.template();
		this.container.appendChild(this.element);

		this.instance = BootstrapToast.getOrCreateInstance(this.element, {
			delay: 3000,
			...bootstrapOptions,
		});

		if (autoShow) this.instance.show();

		this.element.addEventListener('hidden.bs.toast', () => {
			this.destroy();
		});
	}

	destroy() {
		this.instance.dispose();
		this.element.remove();
	}

	show() {
		this.instance.show();
	}

	hide() {
		this.instance.hide();
	}

	getVariantIcon(): IconProps['icon'] {
		switch (this.variant) {
			case 'info':
				return 'info-circle';
			case 'success':
				return 'check-circle';
			case 'error':
				return 'exclamation-circle';
			case 'warning':
				return 'triangle-exclamation';
		}
	}

	template() {
		return (
			<div
				className={clsx('toast position-relative bottom-0 end-0', `toast--${this.variant}`, this.additionalClasses)}
				attributes={{
					role: 'alert',
					'aria-live': 'assertive',
					'aria-atomic': 'true',
				}}>
				<div className="toast-header">
					<Icon icon={this.getVariantIcon()} size="2xl" className="d-block me-2" />
					<strong className="me-auto">{this.title}</strong>
					{this.canClose && (
						<Button
							variant="close"
							attributes={{
								'aria-label': 'Close',
							}}
							dataset={{
								bsDismiss: 'toast',
							}}
							iconLeft={<Icon icon={this.getVariantIcon()} size="xl" />}
						/>
					)}
				</div>
				<div className="toast-body">{this.body}</div>
			</div>
		) as HTMLElement;
	}
}

export default Toast;
